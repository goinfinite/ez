package backupInfra

import (
	"encoding/json"
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

func (repo *BackupCmdRepo) userDataWatermarkLimitValidator() error {
	userDataDirectoryWatermarkLimitPercent := float64(92)
	userDataDirectoryStats, err := repo.readUserDataStats()
	if err != nil {
		return err
	}

	if userDataDirectoryStats.UsedPercent >= userDataDirectoryWatermarkLimitPercent {
		return errors.New("UserDataDirectoryUsageExceedsWatermarkLimit")
	}

	return nil
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
) (archiveFile entity.ContainerImageArchive, err error) {
	shouldCreateArchive := false
	createSnapshotDto := dto.CreateContainerSnapshotImage{
		ContainerId:         containerWithMetrics.Id,
		ShouldCreateArchive: &shouldCreateArchive,
	}
	snapshotImageId, err := containerImageCmdRepo.CreateSnapshot(createSnapshotDto)
	if err != nil {
		return archiveFile, errors.New("CreateSnapshotImageFailed: " + err.Error())
	}

	createArchiveDto := dto.CreateContainerImageArchive{
		AccountId:       containerWithMetrics.AccountId,
		ImageId:         snapshotImageId,
		DestinationPath: &jobTmpDir,
	}
	archiveFile, err = containerImageCmdRepo.CreateArchive(createArchiveDto)
	if err != nil {
		return archiveFile, errors.New("CreateArchiveFailed: " + err.Error())
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

func (repo *BackupCmdRepo) readTaskRemotePath(
	destinationEntity entity.IBackupDestination,
	taskId valueObject.BackupTaskId,
) string {
	remotePathStr := "encdest:"
	switch destinationEntity := destinationEntity.(type) {
	case entity.BackupDestinationRemoteHost:
		if destinationEntity.DestinationPath.String() != "/" {
			remotePathStr += destinationEntity.DestinationPath.String()
		}
	}

	return remotePathStr + taskId.String()
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
	archiveEntity entity.ContainerImageArchive,
) BackupTaskRunDetails {
	containerIdStr := containerWithMetrics.Id.String()

	var backupBinaryFlags []string
	switch taskRunDetails.DestinationEntity.(type) {
	case entity.BackupDestinationObjectStorage:
		backupBinaryFlags = append(backupBinaryFlags, "--s3-no-check-bucket")
	}

	archiveFileExtStr := ".tar.br"
	archiveFileExt, err := archiveEntity.UnixFilePath.ReadCompoundFileExtension()
	if err == nil {
		archiveFileExtStr = "." + archiveFileExt.String()
	}

	destFileNameStr := containerWithMetrics.AccountId.String() + "-" + containerIdStr +
		"-" + archiveEntity.ImageId.String() + archiveFileExtStr
	destPathStr := repo.readTaskRemotePath(taskRunDetails.DestinationEntity, taskRunDetails.TaskId)
	destPathStr += "/" + destFileNameStr

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
	backupBinaryCmd += " copyto " + archiveEntity.UnixFilePath.String() + " " + destPathStr

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
	taskRunDetails.SizeBytes += archiveEntity.SizeBytes.Uint64()
	taskRunDetails.FinishedAt = time.Now()
	taskRunDetails.ElapsedSecs = uint64(
		taskRunDetails.FinishedAt.Sub(taskRunDetails.StartedAt).Seconds(),
	)

	return taskRunDetails
}

func (repo *BackupCmdRepo) RunJob(runDto dto.RunBackupJob) error {
	err := repo.userDataWatermarkLimitValidator()
	if err != nil {
		return err
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

			archiveEntity, err := repo.createContainerArchive(
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
					preUploadTaskRunDetails, containerWithMetrics, archiveEntity,
				)
				taskIdRunDetailsMap[taskId] = postUploadTaskRunDetails
			}

			err = containerImageCmdRepo.DeleteArchive(archiveEntity)
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

func (repo *BackupCmdRepo) restoreContainerArchive(
	containerArchiveEntity entity.ContainerImageArchive,
	containerQueryRepo *infra.ContainerQueryRepo,
	containerCmdRepo *infra.ContainerCmdRepo,
	containerImageQueryRepo *infra.ContainerImageQueryRepo,
	containerImageCmdRepo *infra.ContainerImageCmdRepo,
	mappingQueryRepo *infra.MappingQueryRepo,
	mappingCmdRepo *infra.MappingCmdRepo,
	shouldReplaceExistingContainers bool,
	shouldRestoreMappings bool,
) error {
	importImageArchiveDto := dto.ImportContainerImageArchive{
		AccountId:       containerArchiveEntity.AccountId,
		ArchiveFilePath: &containerArchiveEntity.UnixFilePath,
	}
	imageId, err := containerImageCmdRepo.ImportArchive(importImageArchiveDto)
	if err != nil {
		return errors.New("ImportContainerImageArchiveFailed: " + err.Error())
	}

	containerImageEntity, err := containerImageQueryRepo.ReadById(
		containerArchiveEntity.AccountId, imageId,
	)
	if err != nil {
		return errors.New("ReadContainerImageFailed: " + err.Error())
	}

	if containerImageEntity.OriginContainerDetails == nil {
		return errors.New("OriginContainerDetailsNotFound")
	}

	rawContainerHostname := containerImageEntity.OriginContainerDetails.Hostname.String()
	if !shouldReplaceExistingContainers {
		archiveCreatedAtStr := containerArchiveEntity.CreatedAt.String()
		rawContainerHostname = archiveCreatedAtStr + ".restored." + rawContainerHostname
	}
	containerHostname, err := valueObject.NewFqdn(rawContainerHostname)
	if err != nil {
		return errors.New("ValidateContainerHostnameFailed: " + err.Error())
	}

	if shouldReplaceExistingContainers {
		requestContainerDto := dto.ReadContainersRequest{
			ContainerId: []valueObject.ContainerId{
				containerImageEntity.OriginContainerDetails.Id,
			},
		}
		_, err = containerQueryRepo.ReadFirst(requestContainerDto)
		if err == nil {
			deleteContainerDto := dto.DeleteContainer{
				AccountId:   containerImageEntity.OriginContainerDetails.AccountId,
				ContainerId: containerImageEntity.OriginContainerDetails.Id,
			}

			err = containerCmdRepo.Delete(deleteContainerDto)
			if err != nil {
				return errors.New("DeleteContainerFailed: " + err.Error())
			}
		}
	}

	createContainerDto := dto.CreateContainer{
		AccountId:          containerArchiveEntity.AccountId,
		Hostname:           containerHostname,
		ImageId:            &imageId,
		PortBindings:       containerImageEntity.PortBindings,
		Envs:               containerImageEntity.Envs,
		Entrypoint:         containerImageEntity.Entrypoint,
		ProfileId:          &containerImageEntity.OriginContainerDetails.ProfileId,
		RestartPolicy:      &containerImageEntity.OriginContainerDetails.RestartPolicy,
		AutoCreateMappings: false,
	}

	containerId, err := containerCmdRepo.Create(createContainerDto)
	if err != nil {
		return errors.New("RestoreContainerFailed: " + err.Error())
	}

	err = os.Remove(containerArchiveEntity.UnixFilePath.String())
	if err != nil {
		return errors.New("DeleteContainerImageArchiveFailed: " + err.Error())
	}

	if !shouldRestoreMappings {
		return nil
	}

	for _, mappingEntity := range containerImageEntity.OriginContainerMappings {
		currentMappingEntity, err := mappingQueryRepo.ReadById(mappingEntity.Id)
		if err == nil {
			isSameHostname := currentMappingEntity.Hostname == mappingEntity.Hostname
			isSamePublicPort := currentMappingEntity.PublicPort == mappingEntity.PublicPort
			if isSameHostname && isSamePublicPort {
				createMappingTargetDto := dto.CreateMappingTarget{
					AccountId:   containerArchiveEntity.AccountId,
					MappingId:   currentMappingEntity.Id,
					ContainerId: containerId,
				}
				_, err = mappingCmdRepo.CreateTarget(createMappingTargetDto)
				if err != nil {
					slog.Debug(
						"RestoreMappingTargetFailed",
						slog.Uint64("currentMappingId", currentMappingEntity.Id.Uint64()),
						slog.String("containerId", containerId.String()),
						slog.String("error", err.Error()),
					)
				}

				continue
			}
		}

		createMappingDto := dto.CreateMapping{
			AccountId:    containerArchiveEntity.AccountId,
			Hostname:     mappingEntity.Hostname,
			PublicPort:   mappingEntity.PublicPort,
			Protocol:     mappingEntity.Protocol,
			ContainerIds: []valueObject.ContainerId{containerId},
		}
		_, err = mappingCmdRepo.Create(createMappingDto)
		if err != nil {
			slog.Debug(
				"RestoreMappingFailed",
				slog.Uint64("publicPort", uint64(mappingEntity.PublicPort.Uint16())),
				slog.String("containerId", containerId.String()),
				slog.String("error", err.Error()),
			)
			continue
		}
	}

	return nil
}

func (repo *BackupCmdRepo) RestoreTask(restoreDto dto.RestoreBackupTask) error {
	taskArchiveProvided := restoreDto.ArchiveId != nil
	if !taskArchiveProvided {
		createArchiveDto := dto.CreateBackupTaskArchive{
			TaskId:                    *restoreDto.TaskId,
			TimeoutSecs:               restoreDto.TimeoutSecs,
			ContainerAccountIds:       restoreDto.ContainerAccountIds,
			ContainerIds:              restoreDto.ContainerIds,
			ExceptContainerAccountIds: restoreDto.ExceptContainerAccountIds,
			ExceptContainerIds:        restoreDto.ExceptContainerIds,
			OperatorAccountId:         valueObject.SystemAccountId,
		}

		archiveId, err := repo.CreateTaskArchive(createArchiveDto)
		if err != nil {
			return errors.New("CreateTaskArchiveFailed: " + err.Error())
		}
		restoreDto.ArchiveId = &archiveId
	}

	archiveEntity, err := repo.backupQueryRepo.ReadFirstTaskArchive(
		dto.ReadBackupTaskArchivesRequest{ArchiveId: restoreDto.ArchiveId},
	)
	if err != nil {
		return errors.New("ReadBackupTaskArchiveFailed: " + err.Error())
	}

	userDataDirectoryStats, err := repo.readUserDataStats()
	if err != nil {
		return errors.New("ReadUserDataDirStatsError: " + err.Error())
	}
	// @see https://ntorga.com/gzip-bzip2-xz-zstd-7z-brotli-or-lz4/
	necessaryFreeStorageBytes := archiveEntity.SizeBytes.Uint64() * 3
	if necessaryFreeStorageBytes > userDataDirectoryStats.Free {
		return errors.New("InsufficientUserDataDirectoryFreeSpace")
	}

	restoreBaseTmpDir, err := valueObject.NewUnixFilePath(infraEnvs.RestoreTaskTmpDir)
	if err != nil {
		return errors.New("ValidateRestoreBaseTaskTmpDirFailed: " + err.Error())
	}
	restoreBaseTmpDirStr := restoreBaseTmpDir.String()
	err = infraHelper.MakeDir(restoreBaseTmpDirStr)
	if err != nil {
		return errors.New("MakeRestoreBaseTaskTmpDirFailed: " + err.Error())
	}

	_, err = infraHelper.RunCmd("chown", "nobody:nogroup", restoreBaseTmpDirStr)
	if err != nil {
		return errors.New("ChownRestoreBaseTaskTmpDirFailed: " + err.Error())
	}

	_, err = infraHelper.RunCmd(
		"tar", "-xf", archiveEntity.UnixFilePath.String(), "-C", restoreBaseTmpDirStr,
	)
	if err != nil {
		return errors.New("ExtractTaskArchiveFailed: " + err.Error())
	}

	if !taskArchiveProvided {
		err = os.Remove(archiveEntity.UnixFilePath.String())
		if err != nil {
			return errors.New("DeleteTempTaskArchiveFailed: " + err.Error())
		}
	}

	rawRestoreTmpDir := restoreBaseTmpDirStr + "/" +
		archiveEntity.UnixFilePath.ReadFileNameWithoutExtension().String()
	restoreTaskTmpDir, err := valueObject.NewUnixFilePath(rawRestoreTmpDir)
	if err != nil {
		return errors.New("ValidateRestoreTmpDirFailed: " + err.Error())
	}

	operatorAccountIdInt := int(restoreDto.OperatorAccountId)
	err = os.Chown(
		restoreTaskTmpDir.String(), operatorAccountIdInt, operatorAccountIdInt,
	)
	if err != nil {
		return errors.New("ChownRestoreTmpDirFailed: " + err.Error())
	}

	containerImageQueryRepo := infra.NewContainerImageQueryRepo(repo.persistentDbSvc)
	requestContainerImagesDto := dto.ReadContainerImageArchivesRequest{
		Pagination: dto.Pagination{
			ItemsPerPage: 1000,
		},
		ArchivesDirectory: &restoreTaskTmpDir,
	}

	containerImagesResponseDto, err := containerImageQueryRepo.ReadArchives(
		requestContainerImagesDto,
	)
	if err != nil {
		return errors.New("ReadContainerImageArchivesFailed: " + err.Error())
	}

	if len(containerImagesResponseDto.Archives) == 0 {
		return errors.New("NoContainerImageArchivesFound")
	}

	shouldReplaceExistingContainers := false
	if restoreDto.ShouldReplaceExistingContainers != nil {
		shouldReplaceExistingContainers = *restoreDto.ShouldReplaceExistingContainers
	}

	shouldRestoreMappings := true
	if restoreDto.ShouldRestoreMappings != nil {
		shouldRestoreMappings = *restoreDto.ShouldRestoreMappings
	}

	containerQueryRepo := infra.NewContainerQueryRepo(repo.persistentDbSvc)
	containerCmdRepo := infra.NewContainerCmdRepo(repo.persistentDbSvc)
	containerImageCmdRepo := infra.NewContainerImageCmdRepo(repo.persistentDbSvc)
	mappingQueryRepo := infra.NewMappingQueryRepo(repo.persistentDbSvc)
	mappingCmdRepo := infra.NewMappingCmdRepo(repo.persistentDbSvc)

	for _, imageArchiveEntity := range containerImagesResponseDto.Archives {
		err = repo.restoreContainerArchive(
			imageArchiveEntity, containerQueryRepo, containerCmdRepo,
			containerImageQueryRepo, containerImageCmdRepo,
			mappingQueryRepo, mappingCmdRepo,
			shouldReplaceExistingContainers, shouldRestoreMappings,
		)
		if err != nil {
			slog.Debug(
				"RestoreContainerArchiveFailed",
				slog.String("imageArchiveEntityId", imageArchiveEntity.ImageId.String()),
				slog.String("error", err.Error()),
			)
			continue
		}
	}

	err = os.RemoveAll(restoreTaskTmpDir.String())
	if err != nil {
		return errors.New("RemoveRestoreTmpDirFailed: " + err.Error())
	}

	return nil
}

func (repo *BackupCmdRepo) DeleteTask(deleteDto dto.DeleteBackupTask) error {
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

	remotePathStr := repo.readTaskRemotePath(destinationEntity, taskEntity.TaskId)
	backupBinaryCmd := backupBinaryCli + " delete " + remotePathStr
	_, err = infraHelper.RunCmdAsUserWithSubShell(taskEntity.AccountId, backupBinaryCmd)
	if err != nil {
		return errors.New("DeleteBackupTaskFilesFailed: " + err.Error())
	}

	return nil
}

func (repo *BackupCmdRepo) readRemoteContainerArchives(
	taskEntity entity.BackupTask,
	destinationEntity entity.IBackupDestination,
	containerAccountIds []valueObject.AccountId,
	containerIds []valueObject.ContainerId,
	exceptContainerAccountIds []valueObject.AccountId,
	exceptContainerIds []valueObject.ContainerId,
) ([]entity.ContainerImageArchive, error) {
	containerArchives := []entity.ContainerImageArchive{}

	taskRemotePath := repo.readTaskRemotePath(destinationEntity, taskEntity.TaskId)

	backupBinaryCli, err := repo.backupBinaryCliFactory(destinationEntity)
	if err != nil {
		return containerArchives, errors.New("BackupCliFactoryFailed: " + err.Error())
	}

	backupBinaryCmd := backupBinaryCli + " lsjson " +
		"--files-only --no-mimetype --no-modtime " + taskRemotePath

	// [
	// {"Path":"1000-677403a6cade-4b2a68046d58.tar.br","Name":"1000-677403a6cade-4b2a68046d58.tar.br","Size":407824652,"ModTime":"","IsDir":false,"Tier":"STANDARD"},
	// {"Path":"1000-b50f1fda0c69-0db3107e7d2c.tar.br","Name":"1000-b50f1fda0c69-0db3107e7d2c.tar.br","Size":252874653,"ModTime":"","IsDir":false,"Tier":"STANDARD"}
	// ]
	rawArchiveFilesList, err := infraHelper.RunCmdAsUserWithSubShell(
		taskEntity.AccountId, backupBinaryCmd,
	)
	if err != nil || len(rawArchiveFilesList) == 0 {
		return containerArchives, errors.New("ListBackupTaskFilesFailed: " + err.Error())
	}

	var rawArchiveFilesListMap []map[string]interface{}
	err = json.Unmarshal([]byte(rawArchiveFilesList), &rawArchiveFilesListMap)
	if err != nil {
		return containerArchives, errors.New("UnmarshalBackupTaskFilesListFailed: " + err.Error())
	}

	accountIdsFilterProvided := len(containerAccountIds) > 0
	containerIdsFilterProvided := len(containerIds) > 0
	exceptAccountIdsFilterProvided := len(exceptContainerAccountIds) > 0
	exceptContainerIdsFilterProvided := len(exceptContainerIds) > 0

	accountIdsFilterSet := infraHelper.SliceToSetTransformer(containerAccountIds)
	containerIdsFilterSet := infraHelper.SliceToSetTransformer(containerIds)
	exceptAccountIdsFilterSet := infraHelper.SliceToSetTransformer(exceptContainerAccountIds)
	exceptContainerIdsFilterSet := infraHelper.SliceToSetTransformer(exceptContainerIds)

	for _, rawArchiveFile := range rawArchiveFilesListMap {
		rawName, assertOk := rawArchiveFile["Name"].(string)
		if !assertOk {
			slog.Debug("AssertNameFailed", slog.Any("rawArchiveFile", rawArchiveFile))
			continue
		}

		nameParts := strings.Split(rawName, "-")
		if len(nameParts) < 3 {
			slog.Debug("SplitNamePartsFailed", slog.String("rawName", rawName))
			continue
		}

		accountId, err := valueObject.NewAccountId(nameParts[0])
		if err != nil {
			slog.Debug(err.Error(), slog.String("rawAccountId", nameParts[0]))
			continue
		}

		if accountIdsFilterProvided {
			if _, exists := accountIdsFilterSet[accountId]; !exists {
				continue
			}
		}

		if exceptAccountIdsFilterProvided {
			if _, exists := exceptAccountIdsFilterSet[accountId]; exists {
				continue
			}
		}

		containerId, err := valueObject.NewContainerId(nameParts[1])
		if err != nil {
			slog.Debug(err.Error(), slog.String("rawContainerId", nameParts[1]))
			continue
		}

		if containerIdsFilterProvided {
			if _, exists := containerIdsFilterSet[containerId]; !exists {
				continue
			}
		}

		if exceptContainerIdsFilterProvided {
			if _, exists := exceptContainerIdsFilterSet[containerId]; exists {
				continue
			}
		}

		sizeBytes, err := valueObject.NewByte(rawArchiveFile["Size"])
		if err != nil {
			slog.Debug(err.Error(), slog.Any("sizeBytes", rawArchiveFile["Size"]))
			continue
		}

		rawPath, assertOk := rawArchiveFile["Path"].(string)
		if !assertOk {
			slog.Debug("AssertPathFailed", slog.Any("rawArchiveFile", rawArchiveFile))
			continue
		}
		if !strings.HasPrefix(rawPath, "/") {
			rawPath = "/" + rawPath
		}
		filePath, err := valueObject.NewUnixFilePath(rawPath)
		if err != nil {
			slog.Debug(err.Error(), slog.String("rawPath", rawPath))
			continue
		}

		imageArchiveFile := entity.ContainerImageArchive{
			AccountId:    accountId,
			UnixFilePath: filePath,
			SizeBytes:    sizeBytes,
			ContainerId:  &containerId,
		}
		containerArchives = append(containerArchives, imageArchiveFile)
	}

	return containerArchives, nil
}

func (repo *BackupCmdRepo) CreateTaskArchive(
	createDto dto.CreateBackupTaskArchive,
) (archiveId valueObject.BackupTaskArchiveId, err error) {
	taskEntity, err := repo.backupQueryRepo.ReadFirstTask(
		dto.ReadBackupTasksRequest{TaskId: &createDto.TaskId},
	)
	if err != nil {
		return archiveId, errors.New("ReadBackupTaskFailed: " + err.Error())
	}

	userDataDirectoryStats, err := repo.readUserDataStats()
	if err != nil {
		return archiveId, errors.New("ReadUserDataDirStatsError: " + err.Error())
	}

	accountIdsFilterProvided := len(createDto.ContainerAccountIds) > 0
	containerIdsFilterProvided := len(createDto.ContainerIds) > 0
	exceptAccountIdsFilterProvided := len(createDto.ExceptContainerAccountIds) > 0
	exceptContainerIdsFilterProvided := len(createDto.ExceptContainerIds) > 0

	wereFiltersProvided := accountIdsFilterProvided || containerIdsFilterProvided ||
		exceptAccountIdsFilterProvided || exceptContainerIdsFilterProvided
	necessaryFreeStorageBytes := taskEntity.SizeBytes.Uint64() * 2
	if !wereFiltersProvided && necessaryFreeStorageBytes >= userDataDirectoryStats.Free {
		return archiveId, errors.New("InsufficientUserDataDirectoryFreeSpace")
	}

	destinationEntity, err := repo.backupQueryRepo.ReadFirstDestination(
		dto.ReadBackupDestinationsRequest{DestinationId: &taskEntity.DestinationId}, true,
	)
	if err != nil {
		return archiveId, errors.New("ReadBackupDestinationFailed: " + err.Error())
	}

	necessaryFreeStorageBytes = 0
	containerArchiveEntities, err := repo.readRemoteContainerArchives(
		taskEntity, destinationEntity, createDto.ContainerAccountIds, createDto.ContainerIds,
		createDto.ExceptContainerAccountIds, createDto.ExceptContainerIds,
	)
	if err != nil {
		return archiveId, errors.New("ReadRemoteContainerArchivesFailed: " + err.Error())
	}

	if len(containerArchiveEntities) == 0 {
		return archiveId, errors.New("NoContainerArchivesFound")
	}

	for _, containerArchiveEntity := range containerArchiveEntities {
		necessaryFreeStorageBytes += containerArchiveEntity.SizeBytes.Uint64()
	}
	if necessaryFreeStorageBytes*2 >= userDataDirectoryStats.Free {
		return archiveId, errors.New("InsufficientUserDataDirectoryFreeSpace")
	}

	archivesDirectoryStr := infraEnvs.NobodyDataDirectory + "/archives"
	if createDto.OperatorAccountId != valueObject.SystemAccountId {
		accountQueryRepo := infra.NewAccountQueryRepo(repo.persistentDbSvc)
		operatorAccountEntity, err := accountQueryRepo.ReadById(createDto.OperatorAccountId)
		if err != nil {
			return archiveId, errors.New("ReadOperatorAccountFailed: " + err.Error())
		}
		if necessaryFreeStorageBytes*2 >= operatorAccountEntity.Quota.StorageBytes.Uint64() {
			return archiveId, errors.New("InsufficientOperatorAccountQuota")
		}

		archivesDirectoryStr = operatorAccountEntity.HomeDirectory.String() + "/archives"
	}
	archivesDirectoryStr += "/tasks"

	createDtoJson, err := json.Marshal(createDto)
	if err != nil {
		return archiveId, errors.New("MarshalCreateDtoFailed: " + err.Error())
	}
	rawOperationHash := infraHelper.GenStrongShortHash(string(createDtoJson))
	archiveId, err = valueObject.NewBackupTaskArchiveId(rawOperationHash)
	if err != nil {
		return archiveId, errors.New("ValidateArchiveTaskIdFailed: " + err.Error())
	}
	archiveIdStr := archiveId.String()
	taskIdStr := taskEntity.TaskId.String()

	taskArchiveDirName := taskIdStr + "-" + archiveIdStr
	rawArchivesDirTaskDir := archivesDirectoryStr + "/" + taskArchiveDirName
	archivesDirTaskDir, err := valueObject.NewUnixFilePath(rawArchivesDirTaskDir)
	if err != nil {
		return archiveId, errors.New("ValidateArchiveTaskDirPathFailed: " + err.Error())
	}
	archivesDirTaskDirStr := archivesDirTaskDir.String()

	err = infraHelper.MakeDir(archivesDirTaskDirStr)
	if err != nil {
		return archiveId, errors.New("MakeArchiveTaskDirFailed: " + err.Error())
	}

	backupBinaryCli, err := repo.backupBinaryCliFactory(destinationEntity)
	if err != nil {
		return archiveId, errors.New("BackupCliFactoryFailed: " + err.Error())
	}

	taskRemotePathStr := repo.readTaskRemotePath(destinationEntity, taskEntity.TaskId)

	for _, containerArchiveEntity := range containerArchiveEntities {
		archiveFilePath := containerArchiveEntity.UnixFilePath

		backupBinaryCmd := backupBinaryCli + " copyto " +
			taskRemotePathStr + archiveFilePath.String() +
			" " + archivesDirTaskDirStr + "/" + archiveFilePath.ReadFileName().String()

		_, err = infraHelper.RunCmdWithSubShell(backupBinaryCmd)
		if err != nil {
			slog.Debug(
				"DownloadContainerArchiveFileFailed",
				slog.String("containerId", containerArchiveEntity.ContainerId.String()),
				slog.String("archiveFilePath", archiveFilePath.String()),
				slog.String("error", err.Error()),
			)
			continue
		}
	}

	rawTaskArchiveFilePath := archivesDirectoryStr + "/" + taskIdStr + "-" + archiveIdStr + ".tar"
	taskArchiveFilePath, err := valueObject.NewUnixFilePath(rawTaskArchiveFilePath)
	if err != nil {
		return archiveId, errors.New("ValidateArchiveTaskFilePathFailed: " + err.Error())
	}

	_, err = infraHelper.RunCmdWithSubShell(
		"tar -cf " + taskArchiveFilePath.String() +
			" -C " + archivesDirectoryStr + " " + taskArchiveDirName,
	)
	if err != nil {
		return archiveId, errors.New("CompressArchiveTaskDirFailed: " + err.Error())
	}

	err = os.RemoveAll(archivesDirTaskDirStr)
	if err != nil {
		return archiveId, errors.New("RemoveArchiveTaskDirFailed: " + err.Error())
	}

	operatorAccountIdInt := int(createDto.OperatorAccountId)
	err = os.Chown(taskArchiveFilePath.String(), operatorAccountIdInt, operatorAccountIdInt)
	if err != nil {
		return archiveId, errors.New("ChownArchiveTaskFileFailed: " + err.Error())
	}

	return archiveId, err
}

func (repo *BackupCmdRepo) DeleteTaskArchive(
	deleteDto dto.DeleteBackupTaskArchive,
) error {
	taskArchiveEntity, err := repo.backupQueryRepo.ReadFirstTaskArchive(
		dto.ReadBackupTaskArchivesRequest{ArchiveId: &deleteDto.ArchiveId},
	)
	if err != nil {
		return errors.New("BackupTaskArchiveNotFound")
	}

	err = os.Remove(taskArchiveEntity.UnixFilePath.String())
	if err != nil {
		return errors.New("RemoveTaskArchiveFileFailed: " + err.Error())
	}

	return nil
}
