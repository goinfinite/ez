package backupInfra

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/user"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/alessio/shellescape"
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/useCase"
	"github.com/goinfinite/ez/src/domain/valueObject"
	"github.com/goinfinite/ez/src/infra"
	"github.com/goinfinite/ez/src/infra/db"
	dbModel "github.com/goinfinite/ez/src/infra/db/model"
	infraEnvs "github.com/goinfinite/ez/src/infra/envs"
	infraHelper "github.com/goinfinite/ez/src/infra/helper"
	"github.com/shirou/gopsutil/disk"
)

type BackupCmdRepo struct {
	persistentDbSvc *db.PersistentDatabaseService
	backupQueryRepo *BackupQueryRepo
}

func NewBackupCmdRepo(
	persistentDbSvc *db.PersistentDatabaseService,
) *BackupCmdRepo {
	return &BackupCmdRepo{
		persistentDbSvc: persistentDbSvc,
		backupQueryRepo: NewBackupQueryRepo(persistentDbSvc),
	}
}

func (repo *BackupCmdRepo) CreateDestination(
	createDto dto.CreateBackupDestinationRequest,
) (responseDto dto.CreateBackupDestinationResponse, err error) {
	var descriptionPtr *string
	if createDto.DestinationDescription != nil {
		description := createDto.DestinationDescription.String()
		descriptionPtr = &description
	}

	destinationPathStr := "/"
	if createDto.DestinationPath != nil {
		destinationPathStr = createDto.DestinationPath.String()
	}

	dbEncryptSecret := os.Getenv("BACKUP_KEYS_SECRET")
	if dbEncryptSecret == "" {
		return responseDto, errors.New("BackupKeysSecretMissing")
	}

	rawEncryptionKey := infraHelper.GenPass(32)
	encryptionKey, err := valueObject.NewPassword(rawEncryptionKey)
	if err != nil {
		return responseDto, errors.New("CreateEncryptionKeyFailed: " + err.Error())
	}

	encryptedEncryptionKeyStr, err := infraHelper.EncryptStr(
		dbEncryptSecret, encryptionKey.String(),
	)
	if err != nil {
		return responseDto, errors.New("EncryptEncryptionKeyFailed: " + err.Error())
	}

	var objectStorageProviderPtr, objectStorageProviderRegionPtr *string
	if createDto.ObjectStorageProvider != nil {
		objectStorageProvider := createDto.ObjectStorageProvider.String()
		objectStorageProviderPtr = &objectStorageProvider
	}
	if createDto.ObjectStorageProviderRegion != nil {
		objectStorageProviderRegion := createDto.ObjectStorageProviderRegion.String()
		objectStorageProviderRegionPtr = &objectStorageProviderRegion
	}

	var objectStorageProviderAccessKeyIdPtr, objectStorageProviderSecretAccessKeyPtr *string
	if createDto.ObjectStorageProviderAccessKeyId != nil {
		objectStorageProviderAccessKeyId := createDto.ObjectStorageProviderAccessKeyId.String()
		objectStorageProviderAccessKeyIdPtr = &objectStorageProviderAccessKeyId
	}

	if createDto.ObjectStorageProviderSecretAccessKey != nil {
		encryptedProviderSecretAccessKey, err := infraHelper.EncryptStr(
			dbEncryptSecret, createDto.ObjectStorageProviderSecretAccessKey.String(),
		)
		if err != nil {
			return responseDto, errors.New("EncryptProviderSecretAccessKeyFailed: " + err.Error())
		}
		objectStorageProviderSecretAccessKeyPtr = &encryptedProviderSecretAccessKey
	}

	var objectStorageEndpointUrlPtr, objectStorageBucketNamePtr *string
	if createDto.ObjectStorageEndpointUrl != nil {
		objectStorageEndpointUrl := createDto.ObjectStorageEndpointUrl.String()
		objectStorageEndpointUrlPtr = &objectStorageEndpointUrl
	}
	if createDto.ObjectStorageBucketName != nil {
		objectStorageBucketName := createDto.ObjectStorageBucketName.String()
		objectStorageBucketNamePtr = &objectStorageBucketName
	}

	var remoteHostTypePtr, remoteHostnamePtr, remoteHostUsernamePtr *string
	if createDto.RemoteHostType != nil {
		remoteHostType := createDto.RemoteHostType.String()
		remoteHostTypePtr = &remoteHostType
	}
	if createDto.RemoteHostname != nil {
		remoteHostname := createDto.RemoteHostname.String()
		remoteHostnamePtr = &remoteHostname
	}
	if createDto.RemoteHostUsername != nil {
		remoteHostUsername := createDto.RemoteHostUsername.String()
		remoteHostUsernamePtr = &remoteHostUsername
	}

	var remoteHostPasswordPtr, remoteHostPrivateKeyFilePathPtr *string
	if createDto.RemoteHostPassword != nil {
		encryptedPassword, err := infraHelper.EncryptStr(
			dbEncryptSecret, createDto.RemoteHostPassword.String(),
		)
		if err != nil {
			return responseDto, errors.New("EncryptPasswordFailed: " + err.Error())
		}
		remoteHostPasswordPtr = &encryptedPassword
	}
	if createDto.RemoteHostPrivateKeyFilePath != nil {
		remoteHostPrivateKeyFilePath := createDto.RemoteHostPrivateKeyFilePath.String()
		remoteHostPrivateKeyFilePathPtr = &remoteHostPrivateKeyFilePath
	}

	var remoteHostNetworkPortPtr *uint16
	if createDto.RemoteHostNetworkPort != nil {
		remoteHostNetworkPort := createDto.RemoteHostNetworkPort.Uint16()
		remoteHostNetworkPortPtr = &remoteHostNetworkPort
	}

	destinationModel := dbModel.NewBackupDestination(
		0, createDto.AccountId.Uint64(), createDto.DestinationName.String(),
		descriptionPtr, createDto.DestinationType.String(), destinationPathStr,
		encryptedEncryptionKeyStr, createDto.MinLocalStorageFreePercent,
		createDto.MaxDestinationStorageUsagePercent, createDto.MaxConcurrentConnections,
		createDto.DownloadBytesSecRateLimit, createDto.UploadBytesSecRateLimit,
		createDto.SkipCertificateVerification, objectStorageProviderPtr,
		objectStorageProviderRegionPtr, objectStorageProviderAccessKeyIdPtr,
		objectStorageProviderSecretAccessKeyPtr, objectStorageEndpointUrlPtr,
		objectStorageBucketNamePtr, remoteHostTypePtr, remoteHostnamePtr,
		remoteHostUsernamePtr, remoteHostPasswordPtr, remoteHostPrivateKeyFilePathPtr,
		remoteHostNetworkPortPtr, createDto.RemoteHostConnectionTimeoutSecs,
		createDto.RemoteHostConnectionRetrySecs,
	)

	err = repo.persistentDbSvc.Handler.Create(&destinationModel).Error
	if err != nil {
		return responseDto, err
	}

	destinationId, err := valueObject.NewBackupDestinationId(destinationModel.ID)
	if err != nil {
		return responseDto, err
	}

	return dto.CreateBackupDestinationResponse{
		DestinationId: destinationId,
		EncryptionKey: encryptionKey,
	}, nil
}

func (repo *BackupCmdRepo) UpdateDestination(
	updateDto dto.UpdateBackupDestination,
) error {
	updateMap := map[string]interface{}{
		"account_id": updateDto.AccountId.Uint64(),
	}

	if updateDto.DestinationName != nil {
		updateMap["name"] = updateDto.DestinationName.String()
	}

	if updateDto.DestinationDescription != nil {
		updateMap["description"] = updateDto.DestinationDescription.String()
	}

	if updateDto.DestinationType != nil {
		updateMap["type"] = updateDto.DestinationType.String()
	}

	if updateDto.DestinationPath != nil {
		updateMap["path"] = updateDto.DestinationPath.String()
	}

	if updateDto.MinLocalStorageFreePercent != nil {
		updateMap["min_local_storage_free_percent"] = *updateDto.MinLocalStorageFreePercent
	}

	if updateDto.MaxDestinationStorageUsagePercent != nil {
		updateMap["max_destination_storage_usage_percent"] = *updateDto.MaxDestinationStorageUsagePercent
	}

	if updateDto.MaxConcurrentConnections != nil {
		updateMap["max_concurrent_connections"] = *updateDto.MaxConcurrentConnections
	}

	if updateDto.TotalSpaceUsageBytes != nil {
		updateMap["total_space_usage_bytes"] = uint64(updateDto.TotalSpaceUsageBytes.Int64())
	}

	if updateDto.TotalSpaceUsagePercent != nil {
		updateMap["total_space_usage_percent"] = *updateDto.TotalSpaceUsagePercent
	}

	if updateDto.DownloadBytesSecRateLimit != nil {
		updateMap["download_bytes_sec_rate_limit"] = *updateDto.DownloadBytesSecRateLimit
	}

	if updateDto.UploadBytesSecRateLimit != nil {
		updateMap["upload_bytes_sec_rate_limit"] = *updateDto.UploadBytesSecRateLimit
	}

	if updateDto.SkipCertificateVerification != nil {
		updateMap["skip_certificate_verification"] = *updateDto.SkipCertificateVerification
	}

	if updateDto.ObjectStorageProvider != nil {
		updateMap["object_storage_provider"] = updateDto.ObjectStorageProvider.String()
	}

	if updateDto.ObjectStorageProviderRegion != nil {
		updateMap["object_storage_provider_region"] = updateDto.ObjectStorageProviderRegion.String()
	}

	if updateDto.ObjectStorageProviderAccessKeyId != nil {
		updateMap["object_storage_provider_access_key_id"] = updateDto.ObjectStorageProviderAccessKeyId.String()
	}

	dbEncryptSecret := os.Getenv("BACKUP_KEYS_SECRET")
	if dbEncryptSecret == "" {
		return errors.New("BackupKeysSecretMissing")
	}

	if updateDto.ObjectStorageProviderSecretAccessKey != nil {
		encryptedProviderSecretAccessKey, err := infraHelper.EncryptStr(
			dbEncryptSecret, updateDto.ObjectStorageProviderSecretAccessKey.String(),
		)
		if err != nil {
			return errors.New("EncryptProviderSecretAccessKeyFailed: " + err.Error())
		}
		updateMap["object_storage_provider_secret_access_key"] = encryptedProviderSecretAccessKey
	}

	if updateDto.ObjectStorageEndpointUrl != nil {
		updateMap["object_storage_endpoint_url"] = updateDto.ObjectStorageEndpointUrl.String()
	}

	if updateDto.ObjectStorageBucketName != nil {
		updateMap["object_storage_bucket_name"] = updateDto.ObjectStorageBucketName.String()
	}

	if updateDto.RemoteHostType != nil {
		updateMap["remote_host_type"] = updateDto.RemoteHostType.String()
	}

	if updateDto.RemoteHostname != nil {
		updateMap["remote_hostname"] = updateDto.RemoteHostname.String()
	}

	if updateDto.RemoteHostNetworkPort != nil {
		updateMap["remote_host_network_port"] = updateDto.RemoteHostNetworkPort.Uint16()
	}

	if updateDto.RemoteHostUsername != nil {
		updateMap["remote_host_username"] = updateDto.RemoteHostUsername.String()
	}

	if updateDto.RemoteHostPassword != nil {
		encryptedPassword, err := infraHelper.EncryptStr(
			dbEncryptSecret, updateDto.RemoteHostPassword.String(),
		)
		if err != nil {
			return errors.New("EncryptPasswordFailed: " + err.Error())
		}
		updateMap["remote_host_password"] = encryptedPassword
	}

	if updateDto.RemoteHostPrivateKeyFilePath != nil {
		updateMap["remote_host_private_key_file_path"] = updateDto.RemoteHostPrivateKeyFilePath.String()
	}

	if updateDto.RemoteHostConnectionTimeoutSecs != nil {
		updateMap["remote_host_connection_timeout_secs"] = *updateDto.RemoteHostConnectionTimeoutSecs
	}

	if updateDto.RemoteHostConnectionRetrySecs != nil {
		updateMap["remote_host_connection_retry_secs"] = *updateDto.RemoteHostConnectionRetrySecs
	}

	return repo.persistentDbSvc.Handler.
		Model(&dbModel.BackupDestination{}).
		Where("id = ?", updateDto.DestinationId.Uint64()).
		Updates(updateMap).Error
}

func (repo *BackupCmdRepo) DeleteDestination(
	deleteDto dto.DeleteBackupDestination,
) error {
	return repo.persistentDbSvc.Handler.Model(&dbModel.BackupDestination{}).Delete(
		"id = ? AND account_id = ?",
		deleteDto.DestinationId.Uint64(), deleteDto.AccountId.Uint64(),
	).Error
}

func (repo *BackupCmdRepo) updateBackupCronFile() error {
	desiredJobStatus := true
	readJobsRequestDto := dto.ReadBackupJobsRequest{
		Pagination: dto.Pagination{
			PageNumber:   0,
			ItemsPerPage: 1000,
		},
		JobStatus: &desiredJobStatus,
	}

	readJobsResponseDto, err := repo.backupQueryRepo.ReadJob(readJobsRequestDto)
	if err != nil {
		return err
	}

	var cronFileContent strings.Builder
	warningMessage := `# WARNING: DO NOT EDIT THIS FILE MANUALLY!
# This file is automatically generated by Infinite Ez.
# Any manual changes will be overwritten.
`
	cronFileContent.WriteString(warningMessage)

	for _, jobEntity := range readJobsResponseDto.Jobs {
		cronLineStr := jobEntity.BackupSchedule.String() + " root " +
			infraEnvs.InfiniteEzBinary + " backup job run " +
			"--account-id " + jobEntity.AccountId.String() + " " +
			"--job-id " + jobEntity.JobId.String() + " >/dev/null 2>&1\n"
		cronFileContent.WriteString(cronLineStr)
	}

	return infraHelper.UpdateFile(
		infraEnvs.BackupCronFilePath, cronFileContent.String(), true,
	)
}

func (repo *BackupCmdRepo) CreateJob(
	createDto dto.CreateBackupJob,
) (backupJobId valueObject.BackupJobId, err error) {
	var jobDescriptionPtr *string
	if createDto.JobDescription != nil {
		jobDescription := createDto.JobDescription.String()
		jobDescriptionPtr = &jobDescription
	}

	archiveCompressionFormat := valueObject.CompressionFormatBrotli
	if createDto.ArchiveCompressionFormat != nil {
		archiveCompressionFormat = *createDto.ArchiveCompressionFormat
	}

	destinationIdsUint64 := []uint64{}
	for _, destinationId := range createDto.DestinationIds {
		destinationIdsUint64 = append(destinationIdsUint64, destinationId.Uint64())
	}

	retentionStrategy := valueObject.BackupRetentionStrategyFull
	if createDto.RetentionStrategy != nil {
		retentionStrategy = *createDto.RetentionStrategy
	}

	timeoutSecs := uint64(48 * 60 * 60)
	if createDto.TimeoutSecs != nil {
		timeoutSecs = *createDto.TimeoutSecs
	}

	var containerAccountIdsUint64 []uint64
	for _, containerAccountId := range createDto.ContainerAccountIds {
		containerAccountIdsUint64 = append(containerAccountIdsUint64, containerAccountId.Uint64())
	}

	var containerIds []string
	for _, containerId := range createDto.ContainerIds {
		containerIds = append(containerIds, containerId.String())
	}

	var exceptContainerAccountIdsUint64 []uint64
	for _, exceptContainerAccountId := range createDto.ExceptContainerAccountIds {
		exceptContainerAccountIdsUint64 = append(exceptContainerAccountIdsUint64, exceptContainerAccountId.Uint64())
	}

	var exceptContainerIds []string
	for _, exceptContainerId := range createDto.ExceptContainerIds {
		exceptContainerIds = append(exceptContainerIds, exceptContainerId.String())
	}

	jobModel := dbModel.NewBackupJob(
		0, createDto.AccountId.Uint64(), true, jobDescriptionPtr, destinationIdsUint64,
		retentionStrategy.String(), createDto.BackupSchedule.String(), archiveCompressionFormat.String(),
		timeoutSecs, createDto.MaxTaskRetentionCount, createDto.MaxTaskRetentionDays,
		createDto.MaxConcurrentCpuCores, containerAccountIdsUint64, containerIds,
		exceptContainerAccountIdsUint64, exceptContainerIds,
	)

	err = repo.persistentDbSvc.Handler.Create(&jobModel).Error
	if err != nil {
		return backupJobId, err
	}

	backupJobId, err = valueObject.NewBackupJobId(jobModel.ID)
	if err != nil {
		return backupJobId, err
	}

	return backupJobId, repo.updateBackupCronFile()
}

func (repo *BackupCmdRepo) UpdateJob(
	updateDto dto.UpdateBackupJob,
) error {
	jobUpdatedModel := dbModel.BackupJob{
		AccountID: updateDto.AccountId.Uint64(),
	}

	if updateDto.JobStatus != nil {
		jobUpdatedModel.JobStatus = *updateDto.JobStatus
	}

	if updateDto.JobDescription != nil {
		jobDescriptionStr := updateDto.JobDescription.String()
		jobUpdatedModel.JobDescription = &jobDescriptionStr
	}

	if updateDto.DestinationIds != nil {
		destinationIdsUint64 := []uint64{}
		for _, destinationId := range updateDto.DestinationIds {
			destinationIdsUint64 = append(destinationIdsUint64, destinationId.Uint64())
		}
		jobUpdatedModel.DestinationIds = destinationIdsUint64
	}

	if updateDto.BackupSchedule != nil {
		backupScheduleStr := updateDto.BackupSchedule.String()
		jobUpdatedModel.BackupSchedule = backupScheduleStr
	}

	if updateDto.TimeoutSecs != nil {
		jobUpdatedModel.TimeoutSecs = *updateDto.TimeoutSecs
	}

	if updateDto.MaxTaskRetentionCount != nil {
		jobUpdatedModel.MaxTaskRetentionCount = updateDto.MaxTaskRetentionCount
	}

	if updateDto.MaxTaskRetentionDays != nil {
		jobUpdatedModel.MaxTaskRetentionDays = updateDto.MaxTaskRetentionDays
	}

	if updateDto.MaxConcurrentCpuCores != nil {
		jobUpdatedModel.MaxConcurrentCpuCores = updateDto.MaxConcurrentCpuCores
	}

	if updateDto.ContainerAccountIds != nil {
		containerAccountIdsUint64 := []uint64{}
		for _, containerAccountId := range updateDto.ContainerAccountIds {
			containerAccountIdsUint64 = append(containerAccountIdsUint64, containerAccountId.Uint64())
		}
		jobUpdatedModel.ContainerAccountIds = containerAccountIdsUint64
	}

	if updateDto.ContainerIds != nil {
		containerIds := []string{}
		for _, containerId := range updateDto.ContainerIds {
			containerIds = append(containerIds, containerId.String())
		}
		jobUpdatedModel.ContainerIds = containerIds
	}

	if updateDto.ExceptContainerAccountIds != nil {
		exceptContainerAccountIdsUint64 := []uint64{}
		for _, exceptContainerAccountId := range updateDto.ExceptContainerAccountIds {
			exceptContainerAccountIdsUint64 = append(exceptContainerAccountIdsUint64, exceptContainerAccountId.Uint64())
		}
		jobUpdatedModel.ExceptContainerAccountIds = exceptContainerAccountIdsUint64
	}

	if updateDto.ExceptContainerIds != nil {
		exceptContainerIds := []string{}
		for _, exceptContainerId := range updateDto.ExceptContainerIds {
			exceptContainerIds = append(exceptContainerIds, exceptContainerId.String())
		}
		jobUpdatedModel.ExceptContainerIds = exceptContainerIds
	}

	err := repo.persistentDbSvc.Handler.
		Model(&dbModel.BackupJob{}).
		Where("id = ?", updateDto.JobId.Uint64()).
		Updates(&jobUpdatedModel).Error
	if err != nil {
		return err
	}

	return repo.updateBackupCronFile()
}

func (repo *BackupCmdRepo) DeleteJob(
	deleteDto dto.DeleteBackupJob,
) error {
	err := repo.persistentDbSvc.Handler.Model(&dbModel.BackupJob{}).Delete(
		"id = ? AND account_id = ?",
		deleteDto.JobId.Uint64(), deleteDto.AccountId.Uint64(),
	).Error
	if err != nil {
		return err
	}

	return repo.updateBackupCronFile()
}

func (repo *BackupCmdRepo) readUserDataStats() (disk.UsageStat, error) {
	userDataDirectoryStats, err := disk.Usage(infraEnvs.UserDataDirectory)
	if err != nil || userDataDirectoryStats == nil {
		return disk.UsageStat{}, errors.New("ReadUserDataDirStatsError: " + err.Error())
	}

	return *userDataDirectoryStats, nil
}

func (repo *BackupCmdRepo) readAccountsContainersWithMetrics(
	jobEntity entity.BackupJob,
) (map[valueObject.AccountId][]dto.ContainerWithMetrics, error) {
	accountIdContainersMap := map[valueObject.AccountId][]dto.ContainerWithMetrics{}
	withMetrics := true
	requestContainersDto := dto.ReadContainersRequest{
		Pagination: dto.Pagination{
			PageNumber:   0,
			ItemsPerPage: 1000,
		},
		ContainerAccountId:       jobEntity.ContainerAccountIds,
		ContainerId:              jobEntity.ContainerIds,
		ExceptContainerAccountId: jobEntity.ExceptContainerAccountIds,
		ExceptContainerId:        jobEntity.ExceptContainerIds,
		WithMetrics:              &withMetrics,
	}
	containerQueryRepo := infra.NewContainerQueryRepo(repo.persistentDbSvc)
	responseContainersDto, err := containerQueryRepo.Read(requestContainersDto)
	if err != nil {
		return accountIdContainersMap, errors.New("ReadContainersFailed: " + err.Error())
	}

	if len(responseContainersDto.ContainersWithMetrics) == 0 {
		return accountIdContainersMap, errors.New("NoContainersFound")
	}

	for _, containerWithMetrics := range responseContainersDto.ContainersWithMetrics {
		accountId := containerWithMetrics.AccountId
		accountIdContainersMap[accountId] = append(
			accountIdContainersMap[accountId], containerWithMetrics,
		)
	}

	for _, containerWithMetrics := range accountIdContainersMap {
		sort.SliceStable(containerWithMetrics, func(i, j int) bool {
			return containerWithMetrics[i].Metrics.StorageSpaceBytes < containerWithMetrics[j].Metrics.StorageSpaceBytes
		})
	}

	return accountIdContainersMap, nil
}

type BackupTaskRunDetails struct {
	TaskId                 valueObject.BackupTaskId
	DestinationEntity      entity.IBackupDestination
	ExecutionOutput        string
	SuccessfulContainerIds []string
	FailedContainerIds     []string
	SizeBytes              uint64
	StartedAt              time.Time
	FinishedAt             time.Time
	ElapsedSecs            uint64
}

func (repo *BackupCmdRepo) backupTaskRunDetailsFactory(
	jobEntity entity.BackupJob,
	operatorAccountId valueObject.AccountId,
) map[valueObject.BackupTaskId]BackupTaskRunDetails {
	taskIdRunDetailsMap := map[valueObject.BackupTaskId]BackupTaskRunDetails{}

	if operatorAccountId == valueObject.SystemAccountId {
		operatorAccountId = jobEntity.AccountId
	}

	for _, destinationId := range jobEntity.DestinationIds {
		startedAt := time.Now()
		taskModel := dbModel.BackupTask{
			AccountID:         operatorAccountId.Uint64(),
			JobID:             jobEntity.JobId.Uint64(),
			DestinationID:     destinationId.Uint64(),
			TaskStatus:        valueObject.BackupTaskStatusExecuting.String(),
			RetentionStrategy: jobEntity.RetentionStrategy.String(),
			BackupSchedule:    jobEntity.BackupSchedule.String(),
			TimeoutSecs:       jobEntity.TimeoutSecs,
			StartedAt:         &startedAt,
		}
		err := repo.persistentDbSvc.Handler.Create(&taskModel).Error
		if err != nil {
			slog.Debug(
				"CreateBackupTaskFailed",
				slog.Uint64("destinationId", destinationId.Uint64()),
				slog.String("error", err.Error()),
			)
			continue
		}

		taskId, err := valueObject.NewBackupTaskId(taskModel.ID)
		if err != nil {
			slog.Debug(err.Error(), slog.Uint64("taskId", taskModel.ID))
			continue
		}

		requestDestinationDto := dto.ReadBackupDestinationsRequest{
			Pagination:    useCase.BackupDestinationsDefaultPagination,
			DestinationId: &destinationId,
		}
		destinationEntity, err := repo.backupQueryRepo.ReadFirstDestination(
			requestDestinationDto, true,
		)
		if err != nil {
			slog.Debug(err.Error(), slog.Uint64("destinationId", destinationId.Uint64()))
			continue
		}

		taskIdRunDetailsMap[taskId] = BackupTaskRunDetails{
			TaskId:                 taskId,
			DestinationEntity:      destinationEntity,
			ExecutionOutput:        "",
			SuccessfulContainerIds: []string{},
			FailedContainerIds:     []string{},
			SizeBytes:              0,
			StartedAt:              startedAt,
			FinishedAt:             startedAt,
			ElapsedSecs:            0,
		}
	}

	return taskIdRunDetailsMap
}

func (repo *BackupCmdRepo) accountStorageAllocator(
	containerWithMetrics dto.ContainerWithMetrics,
	accountEntity entity.Account,
	accountCmdRepo *infra.AccountCmdRepo,
) error {
	userDataDirectoryStats, err := repo.readUserDataStats()
	if err != nil {
		return err
	}

	containerSizeBytesUint := containerWithMetrics.Metrics.StorageSpaceBytes.Uint64()
	rawContainerSizeWithMargin := containerSizeBytesUint + containerSizeBytesUint/10
	containerSizeWithMarginBytes, err := valueObject.NewByte(rawContainerSizeWithMargin)
	if err != nil {
		return errors.New("ContainerSizeWithMarginBytesCreateError")
	}

	if containerSizeWithMarginBytes.Uint64() > userDataDirectoryStats.Free {
		return errors.New("InsufficientUserDataDirectoryFreeSpace")
	}

	accountFreeStorageBytes := accountEntity.Quota.StorageBytes - accountEntity.QuotaUsage.StorageBytes
	if containerSizeBytesUint > accountFreeStorageBytes.Uint64() {
		tempUpdatedStorageBytesQuota := accountEntity.Quota.StorageBytes + containerSizeWithMarginBytes
		tempUpdatedAccountQuota := valueObject.AccountQuota{
			StorageBytes: tempUpdatedStorageBytesQuota,
		}

		err = accountCmdRepo.UpdateQuota(containerWithMetrics.AccountId, tempUpdatedAccountQuota)
		if err != nil {
			return errors.New("UpdateAccountQuotaFailed: " + err.Error())
		}
	}

	return nil
}

func (repo *BackupCmdRepo) createJobTmpDir(tmpDir valueObject.UnixFilePath) error {
	tmpDirStr := tmpDir.String()
	err := infraHelper.MakeDir(tmpDirStr)
	if err != nil {
		return errors.New("MakeBackupJobTmpDirFailed: " + err.Error())
	}

	nobodyUser, err := user.Lookup("nobody")
	if err != nil {
		return errors.New("LookupNobodyUserFailed: " + err.Error())
	}
	nobodyUid, err := strconv.Atoi(nobodyUser.Uid)
	if err != nil {
		return errors.New("ConvertNobodyUidFailed: " + err.Error())
	}

	nogroupGroup, err := user.LookupGroup("nogroup")
	if err != nil {
		return errors.New("LookupNoGroupFailed: " + err.Error())
	}
	nogroupGid, err := strconv.Atoi(nogroupGroup.Gid)
	if err != nil {
		return errors.New("ConvertNoGroupGidFailed: " + err.Error())
	}

	err = os.Chown(tmpDirStr, nobodyUid, nogroupGid)
	if err != nil {
		return errors.New("ChownJobTmpDirFailed: " + err.Error())
	}

	err = os.Chmod(tmpDirStr, 0777)
	if err != nil {
		return errors.New("ChmodJobTmpDirFailed: " + err.Error())
	}

	return nil
}

func (repo *BackupCmdRepo) sharedTaskFailRegister(
	accountId *valueObject.AccountId,
	containerId *valueObject.ContainerId,
	executionOutputPtr *string,
	failedContainerIdsPtr *[]string,
	error error,
) {
	tagIdentifier := ""
	slogAttrs := []interface{}{}

	if accountId != nil {
		accountIdStr := accountId.String()
		tagIdentifier = "[" + accountIdStr + "] "
		slogAttrs = append(slogAttrs, "accountId", accountIdStr)
	}

	if containerId != nil {
		containerIdStr := containerId.String()
		tagIdentifier += "[" + containerIdStr + "] "
		*failedContainerIdsPtr = append(*failedContainerIdsPtr, containerIdStr)
		slogAttrs = append(slogAttrs, "containerId", containerIdStr)
	}

	*executionOutputPtr += tagIdentifier + error.Error() + "\n"
	slog.Debug(error.Error(), slogAttrs...)
}

func (repo *BackupCmdRepo) createContainerArchive(
	containerWithMetrics dto.ContainerWithMetrics,
	containerImageCmdRepo *infra.ContainerImageCmdRepo,
	jobTmpDir valueObject.UnixFilePath,
) (archiveFile entity.ContainerImageArchiveFile, err error) {
	shouldCreateArchive := false
	createSnapshotDto := dto.CreateContainerSnapshotImage{
		ContainerId:         containerWithMetrics.Id,
		ShouldCreateArchive: &shouldCreateArchive,
	}
	snapshotImageId, err := containerImageCmdRepo.CreateSnapshot(createSnapshotDto)
	if err != nil {
		return archiveFile, errors.New("CreateSnapshotImageFailed: " + err.Error())
	}

	createArchiveDto := dto.CreateContainerImageArchiveFile{
		AccountId:       containerWithMetrics.AccountId,
		ImageId:         snapshotImageId,
		DestinationPath: &jobTmpDir,
	}
	archiveFile, err = containerImageCmdRepo.CreateArchiveFile(createArchiveDto)
	if err != nil {
		return archiveFile, errors.New("CreateArchiveFileFailed: " + err.Error())
	}

	deleteSnapshotDto := dto.DeleteContainerImage{
		AccountId: containerWithMetrics.AccountId,
		ImageId:   snapshotImageId,
	}
	err = containerImageCmdRepo.Delete(deleteSnapshotDto)
	if err != nil {
		return archiveFile, errors.New("DeleteSnapshotImageFailed: " + err.Error())
	}

	return archiveFile, nil
}

func (repo *BackupCmdRepo) obscurePassword(
	plainPass valueObject.Password,
) (obscuredPass valueObject.Password, err error) {
	rawObscuredPassword, err := infraHelper.RunCmdWithSubShell(
		"echo '" + plainPass.String() + "' | rclone obscure -",
	)
	if err != nil || len(rawObscuredPassword) == 0 {
		return obscuredPass, errors.New("CreateObscurePasswordFailed")
	}
	rawObscuredPassword = strings.TrimSpace(rawObscuredPassword)

	obscuredPassword, err := valueObject.NewPassword(rawObscuredPassword)
	if err != nil {
		return obscuredPass, errors.New("ValidateObscurePasswordFailed")
	}

	return obscuredPassword, nil
}

func (repo *BackupCmdRepo) taskRunDetailsFailRegister(
	taskRunDetails *BackupTaskRunDetails,
	containerIdStr string,
	errorMessage string,
) {
	taskRunDetails.FailedContainerIds = append(
		taskRunDetails.FailedContainerIds, containerIdStr,
	)
	taskRunDetails.ExecutionOutput += "[" + containerIdStr + "] " + errorMessage + "\n"
	taskRunDetails.FinishedAt = time.Now()
	taskRunDetails.ElapsedSecs = uint64(
		taskRunDetails.FinishedAt.Sub(taskRunDetails.StartedAt).Seconds(),
	)
}

func (repo *BackupCmdRepo) backupBinaryCliFactory(
	destinationIEntity entity.IBackupDestination,
) (string, error) {
	unencryptedDestEnvPrefix := "RCLONE_CONFIG_RAWDEST"
	encryptedDestEnvPrefix := "RCLONE_CONFIG_ENCDEST"

	var backupBinaryEnvs []string

	var encryptionKey valueObject.Password
	switch destinationEntity := destinationIEntity.(type) {
	case entity.BackupDestinationLocal:
		backupBinaryEnvs = []string{
			unencryptedDestEnvPrefix + "_TYPE=local",
		}
		encryptionKey = destinationEntity.EncryptionKey
	case entity.BackupDestinationRemoteHost:
		remoteHostTypeStr := "sftp"
		if destinationEntity.RemoteHostType != nil {
			remoteHostTypeStr = destinationEntity.RemoteHostType.String()
		}

		backupBinaryEnvs = []string{
			unencryptedDestEnvPrefix + "_TYPE=" + remoteHostTypeStr,
			unencryptedDestEnvPrefix + "_HOST=" +
				shellescape.Quote(destinationEntity.RemoteHostname.String()),
			unencryptedDestEnvPrefix + "_USER=" +
				shellescape.Quote(destinationEntity.RemoteHostUsername.String()),
		}
		encryptionKey = destinationEntity.EncryptionKey
		if destinationEntity.RemoteHostNetworkPort != nil {
			backupBinaryEnvs = append(
				backupBinaryEnvs,
				unencryptedDestEnvPrefix+"_PORT="+destinationEntity.RemoteHostNetworkPort.String(),
			)
		}
		if destinationEntity.RemoteHostPassword != nil {
			obscuredPassword, err := repo.obscurePassword(*destinationEntity.RemoteHostPassword)
			if err != nil {
				return "", err
			}
			backupBinaryEnvs = append(
				backupBinaryEnvs, unencryptedDestEnvPrefix+"_PASS="+obscuredPassword.String(),
			)
		}
		if destinationEntity.RemoteHostPrivateKeyFilePath != nil {
			backupBinaryEnvs = append(
				backupBinaryEnvs,
				unencryptedDestEnvPrefix+"_KEY_FILE="+
					shellescape.Quote(destinationEntity.RemoteHostPrivateKeyFilePath.String()),
			)
		}

	case entity.BackupDestinationObjectStorage:
		backupBinaryEnvs = []string{
			unencryptedDestEnvPrefix + "_TYPE=s3",
			unencryptedDestEnvPrefix + "_PROVIDER=Custom",
			unencryptedDestEnvPrefix + "_ACCESS_KEY_ID=" +
				destinationEntity.ObjectStorageProviderAccessKeyId.String(),
			unencryptedDestEnvPrefix + "_SECRET_ACCESS_KEY=" +
				destinationEntity.ObjectStorageProviderSecretAccessKey.String(),
			unencryptedDestEnvPrefix + "_ENDPOINT=" +
				destinationEntity.ObjectStorageEndpointUrl.WithoutSchema(),
		}
		encryptionKey = destinationEntity.EncryptionKey
	default:
		return "", errors.New("InvalidBackupDestinationType")
	}

	encryptionKeyObscured, err := repo.obscurePassword(encryptionKey)
	if err != nil {
		return "", err
	}

	backupBinaryEnvs = append(
		backupBinaryEnvs,
		encryptedDestEnvPrefix+"_TYPE=crypt",
		encryptedDestEnvPrefix+"_DIRECTORY_NAME_ENCRYPTION=false",
		encryptedDestEnvPrefix+"_FILENAME_ENCRYPTION=off",
		encryptedDestEnvPrefix+"_PASSWORD="+encryptionKeyObscured.String(),
	)
	unencryptedDestPathStr := ""
	switch destinationEntity := destinationIEntity.(type) {
	case entity.BackupDestinationLocal:
		unencryptedDestPathStr = destinationEntity.DestinationPath.String()
	case entity.BackupDestinationObjectStorage:
		unencryptedDestPathStr += destinationEntity.ObjectStorageBucketName.String() +
			destinationEntity.DestinationPath.String()
	case entity.BackupDestinationRemoteHost:
		unencryptedDestPathStr = destinationEntity.DestinationPath.String()
		if unencryptedDestPathStr == "/" {
			unencryptedDestPathStr = ""
		}
	}
	backupBinaryEnvs = append(
		backupBinaryEnvs, encryptedDestEnvPrefix+"_REMOTE=rawdest:"+unencryptedDestPathStr,
	)

	backupCliWithEnvs := strings.Join(backupBinaryEnvs, " ") + " rclone"

	return strings.TrimSpace(backupCliWithEnvs), nil
}

func (repo *BackupCmdRepo) uploadContainerArchive(
	taskRunDetails BackupTaskRunDetails,
	containerWithMetrics dto.ContainerWithMetrics,
	archiveFileEntity entity.ContainerImageArchiveFile,
) BackupTaskRunDetails {
	containerIdStr := containerWithMetrics.Id.String()

	var backupBinaryFlags []string
	destPathStr := "encdest:"
	switch destinationEntity := taskRunDetails.DestinationEntity.(type) {
	case entity.BackupDestinationObjectStorage:
		backupBinaryFlags = append(backupBinaryFlags, "--s3-no-check-bucket")
	case entity.BackupDestinationRemoteHost:
		if destinationEntity.DestinationPath.String() != "/" {
			destPathStr += destinationEntity.DestinationPath.String()
		}
	}

	archiveFilePath := archiveFileEntity.UnixFilePath
	srcFilePathStr := archiveFilePath.String()
	destPathStr += taskRunDetails.TaskId.String() + "/" + archiveFilePath.ReadFileName().String()

	backupBinaryCli, err := repo.backupBinaryCliFactory(taskRunDetails.DestinationEntity)
	if err != nil {
		repo.taskRunDetailsFailRegister(
			&taskRunDetails, containerIdStr, "BackupCliFactoryFailed",
		)
		slog.Debug(
			"BackupCliFactoryFailed",
			slog.String("containerId", containerIdStr),
			slog.String("error", err.Error()),
		)
		return taskRunDetails
	}
	backupBinaryCmd := backupBinaryCli
	if len(backupBinaryFlags) > 0 {
		backupBinaryCmd += " " + strings.Join(backupBinaryFlags, " ")
	}
	backupBinaryCmd += " copyto " + srcFilePathStr + " " + destPathStr

	_, err = infraHelper.RunCmdAsUserWithSubShell(
		containerWithMetrics.AccountId, backupBinaryCmd,
	)
	if err != nil {
		repo.taskRunDetailsFailRegister(
			&taskRunDetails, containerIdStr, "UploadArchiveFailed",
		)
		slog.Debug(
			"UploadContainerArchiveFailed",
			slog.String("containerId", containerIdStr),
			slog.String("error", err.Error()),
		)
		return taskRunDetails
	}

	taskRunDetails.SuccessfulContainerIds = append(
		taskRunDetails.SuccessfulContainerIds, containerIdStr,
	)
	taskRunDetails.SizeBytes += archiveFileEntity.SizeBytes.Uint64()
	taskRunDetails.FinishedAt = time.Now()
	taskRunDetails.ElapsedSecs = uint64(
		taskRunDetails.FinishedAt.Sub(taskRunDetails.StartedAt).Seconds(),
	)

	return taskRunDetails
}

func (repo *BackupCmdRepo) RunJob(runDto dto.RunBackupJob) error {
	userDataDirectoryWatermarkLimitPercent := float64(92)
	userDataDirectoryStats, err := repo.readUserDataStats()
	if err != nil {
		return err
	}
	if userDataDirectoryStats.UsedPercent >= userDataDirectoryWatermarkLimitPercent {
		return errors.New("UserDataDirectoryUsageExceedsWatermarkLimit")
	}

	requestJobDto := dto.ReadBackupJobsRequest{
		AccountId: &runDto.AccountId,
		JobId:     &runDto.JobId,
	}
	jobEntity, err := repo.backupQueryRepo.ReadFirstJob(requestJobDto)
	if err != nil {
		return errors.New("ReadBackupJobFailed: " + err.Error())
	}

	accountIdContainerWithMetricsMap, err := repo.readAccountsContainersWithMetrics(jobEntity)
	if err != nil {
		return errors.New("ReadAccountContainersFailed: " + err.Error())
	}

	taskIdRunDetailsMap := repo.backupTaskRunDetailsFactory(
		jobEntity, runDto.OperatorAccountId,
	)
	if len(taskIdRunDetailsMap) == 0 {
		return errors.New("NoBackupTasksCreated")
	}

	rawJobTmpDir := fmt.Sprintf(
		"%s/nobody/backup/%d/%d",
		infraEnvs.UserDataDirectory, runDto.AccountId.Uint64(), runDto.JobId.Uint64(),
	)
	jobTmpDir, err := valueObject.NewUnixFilePath(rawJobTmpDir)
	if err != nil {
		return errors.New("ValidateBackupJobTmpDirFailed: " + err.Error())
	}

	err = repo.createJobTmpDir(jobTmpDir)
	if err != nil {
		return errors.New("CreateBackupJobTmpDirFailed: " + err.Error())
	}

	accountQueryRepo := infra.NewAccountQueryRepo(repo.persistentDbSvc)
	accountCmdRepo := infra.NewAccountCmdRepo(repo.persistentDbSvc)
	containerImageCmdRepo := infra.NewContainerImageCmdRepo(repo.persistentDbSvc)

	sharedExecutionOutputStr := ""
	sharedFailedContainerIdStrs := []string{}

	for accountId, containerWithMetricsSlice := range accountIdContainerWithMetricsMap {
		preTaskAccountEntity, err := accountQueryRepo.ReadById(accountId)
		if err != nil {
			repo.sharedTaskFailRegister(
				&accountId, nil, &sharedExecutionOutputStr, &sharedFailedContainerIdStrs, err,
			)
			continue
		}

		for _, containerWithMetrics := range containerWithMetricsSlice {
			err = repo.accountStorageAllocator(
				containerWithMetrics, preTaskAccountEntity, accountCmdRepo,
			)
			if err != nil {
				repo.sharedTaskFailRegister(
					&containerWithMetrics.AccountId, &containerWithMetrics.Id,
					&sharedExecutionOutputStr, &sharedFailedContainerIdStrs, err,
				)
				continue
			}

			archiveFileEntity, err := repo.createContainerArchive(
				containerWithMetrics, containerImageCmdRepo, jobTmpDir,
			)
			if err != nil {
				repo.sharedTaskFailRegister(
					&containerWithMetrics.AccountId, &containerWithMetrics.Id,
					&sharedExecutionOutputStr, &sharedFailedContainerIdStrs, err,
				)
				continue
			}

			for taskId, preUploadTaskRunDetails := range taskIdRunDetailsMap {
				postUploadTaskRunDetails := repo.uploadContainerArchive(
					preUploadTaskRunDetails, containerWithMetrics, archiveFileEntity,
				)
				taskIdRunDetailsMap[taskId] = postUploadTaskRunDetails
			}

			err = containerImageCmdRepo.DeleteArchiveFile(archiveFileEntity)
			if err != nil {
				repo.sharedTaskFailRegister(
					&containerWithMetrics.AccountId, &containerWithMetrics.Id,
					&sharedExecutionOutputStr, &sharedFailedContainerIdStrs, err,
				)
				continue
			}
		}

		postTaskAccountEntity, err := accountQueryRepo.ReadById(accountId)
		if err != nil {
			repo.sharedTaskFailRegister(
				&accountId, nil, &sharedExecutionOutputStr, &sharedFailedContainerIdStrs, err,
			)
			continue
		}

		if preTaskAccountEntity.Quota.StorageBytes != postTaskAccountEntity.Quota.StorageBytes {
			err = accountCmdRepo.UpdateQuota(accountId, postTaskAccountEntity.Quota)
			if err != nil {
				slog.Debug(
					"RestoreOriginalAccountQuotaFailed",
					slog.Uint64("accountId", accountId.Uint64()),
					slog.String("error", err.Error()),
				)
				continue
			}
		}
	}

	err = os.RemoveAll(jobTmpDir.String())
	if err != nil {
		slog.Debug("RemoveBackupJobTmpDirFailed", slog.String("error", err.Error()))
	}

	for taskId, taskRunDetails := range taskIdRunDetailsMap {
		combinedExecutionOutput := taskRunDetails.ExecutionOutput
		if len(sharedExecutionOutputStr) > 0 {
			combinedExecutionOutput = sharedExecutionOutputStr + "\n" + taskRunDetails.ExecutionOutput
		}
		combinedFailedContainerIds := append(
			sharedFailedContainerIdStrs, taskRunDetails.FailedContainerIds...,
		)

		taskStatus := valueObject.BackupTaskStatusCompleted.String()
		if len(combinedFailedContainerIds) > 0 {
			taskStatus = valueObject.BackupTaskStatusPartial.String()
		}
		if len(taskRunDetails.SuccessfulContainerIds) == 0 {
			taskStatus = valueObject.BackupTaskStatusFailed.String()
		}

		taskModelUpdated := dbModel.BackupTask{
			TaskStatus:             taskStatus,
			ExecutionOutput:        &combinedExecutionOutput,
			SuccessfulContainerIds: taskRunDetails.SuccessfulContainerIds,
			FailedContainerIds:     combinedFailedContainerIds,
			SizeBytes:              &taskRunDetails.SizeBytes,
			ElapsedSecs:            &taskRunDetails.ElapsedSecs,
			FinishedAt:             &taskRunDetails.FinishedAt,
		}

		err := repo.persistentDbSvc.Handler.Model(&dbModel.BackupTask{}).
			Where("id = ?", taskId.Uint64()).
			Updates(&taskModelUpdated).Error
		if err != nil {
			slog.Debug(
				"UpdateBackupTaskModelError: "+err.Error(),
				slog.Uint64("taskId", taskId.Uint64()))
			continue
		}
	}

	return nil
}

func (repo *BackupCmdRepo) DeleteTask(
	deleteDto dto.DeleteBackupTask,
) error {
	taskEntity, err := repo.backupQueryRepo.ReadFirstTask(
		dto.ReadBackupTasksRequest{TaskId: &deleteDto.TaskId},
	)
	if err != nil {
		return errors.New("ReadBackupTaskFailed: " + err.Error())
	}

	err = repo.persistentDbSvc.Handler.Model(&dbModel.BackupTask{}).Delete(
		"id = ?", deleteDto.TaskId.Uint64(),
	).Error
	if err != nil {
		return errors.New("DeleteBackupTaskDatabaseEntryFailed: " + err.Error())
	}

	if !deleteDto.ShouldDiscardFiles {
		return nil
	}

	destinationEntity, err := repo.backupQueryRepo.ReadFirstDestination(
		dto.ReadBackupDestinationsRequest{DestinationId: &taskEntity.DestinationId}, true,
	)
	if err != nil {
		return errors.New("ReadBackupDestinationFailed: " + err.Error())
	}

	backupBinaryCli, err := repo.backupBinaryCliFactory(destinationEntity)
	if err != nil {
		return errors.New("BackupCliFactoryFailed: " + err.Error())
	}

	remotePathStr := "encdest:"
	switch destinationEntity := destinationEntity.(type) {
	case entity.BackupDestinationRemoteHost:
		if destinationEntity.DestinationPath.String() != "/" {
			remotePathStr += destinationEntity.DestinationPath.String()
		}
	}

	backupBinaryCmd := backupBinaryCli + " delete " +
		remotePathStr + taskEntity.TaskId.String() + "/"
	_, err = infraHelper.RunCmdAsUserWithSubShell(taskEntity.AccountId, backupBinaryCmd)
	if err != nil {
		return errors.New("DeleteBackupTaskFilesFailed: " + err.Error())
	}

	return nil
}

func (repo *BackupCmdRepo) CreateTaskArchive(
	createDto dto.CreateBackupTaskArchive,
) (archiveId valueObject.BackupTaskArchiveId, err error) {
	return archiveId, err
}
