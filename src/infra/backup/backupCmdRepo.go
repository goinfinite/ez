package backupInfra

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/user"
	"regexp"
	"slices"
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

const TasksArchivesRelativePath = "/backup/tasks/archives"

type BackupCmdRepo struct {
	persistentDbSvc *db.PersistentDatabaseService
	trailDbSvc      *db.TrailDatabaseService
	backupQueryRepo *BackupQueryRepo
}

func NewBackupCmdRepo(
	persistentDbSvc *db.PersistentDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) *BackupCmdRepo {
	return &BackupCmdRepo{
		persistentDbSvc: persistentDbSvc,
		trailDbSvc:      trailDbSvc,
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

	var downloadBytesSecRateLimitPtr, uploadBytesSecRateLimitPtr *uint64
	if createDto.DownloadBytesSecRateLimit != nil {
		downloadBytesSecRateLimit := createDto.DownloadBytesSecRateLimit.Uint64()
		downloadBytesSecRateLimitPtr = &downloadBytesSecRateLimit
	}
	if createDto.UploadBytesSecRateLimit != nil {
		uploadBytesSecRateLimit := createDto.UploadBytesSecRateLimit.Uint64()
		uploadBytesSecRateLimitPtr = &uploadBytesSecRateLimit
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

	var remoteHostConnectionTimeoutSecsPtr *uint64
	if createDto.RemoteHostConnectionTimeoutSecs != nil {
		remoteHostConnectionTimeoutSecs := createDto.RemoteHostConnectionTimeoutSecs.Uint64()
		remoteHostConnectionTimeoutSecsPtr = &remoteHostConnectionTimeoutSecs
	}

	var remoteHostConnectionRetrySecsPtr *uint64
	if createDto.RemoteHostConnectionRetrySecs != nil {
		remoteHostConnectionRetrySecs := createDto.RemoteHostConnectionRetrySecs.Uint64()
		remoteHostConnectionRetrySecsPtr = &remoteHostConnectionRetrySecs
	}

	destinationModel := dbModel.NewBackupDestination(
		0, createDto.AccountId.Uint64(), createDto.DestinationName.String(),
		descriptionPtr, createDto.DestinationType.String(), destinationPathStr,
		encryptedEncryptionKeyStr, createDto.MinLocalStorageFreePercent,
		createDto.MaxDestinationStorageUsagePercent, createDto.MaxConcurrentConnections,
		downloadBytesSecRateLimitPtr, uploadBytesSecRateLimitPtr,
		createDto.SkipCertificateVerification, objectStorageProviderPtr,
		objectStorageProviderRegionPtr, objectStorageProviderAccessKeyIdPtr,
		objectStorageProviderSecretAccessKeyPtr, objectStorageEndpointUrlPtr,
		objectStorageBucketNamePtr, remoteHostTypePtr, remoteHostnamePtr,
		remoteHostUsernamePtr, remoteHostPasswordPtr, remoteHostPrivateKeyFilePathPtr,
		remoteHostNetworkPortPtr, remoteHostConnectionTimeoutSecsPtr,
		remoteHostConnectionRetrySecsPtr,
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
	destinationUpdatedModel := dbModel.BackupDestination{
		AccountID: updateDto.AccountId.Uint64(),
	}

	if updateDto.DestinationName != nil {
		destinationUpdatedModel.Name = updateDto.DestinationName.String()
	}

	if updateDto.DestinationDescription != nil {
		destinationDescriptionStr := updateDto.DestinationDescription.String()
		destinationUpdatedModel.Description = &destinationDescriptionStr
	}

	if updateDto.DestinationPath != nil {
		destinationUpdatedModel.Path = updateDto.DestinationPath.String()
	}

	if updateDto.MinLocalStorageFreePercent != nil {
		destinationUpdatedModel.MinLocalStorageFreePercent = updateDto.MinLocalStorageFreePercent
	}

	if updateDto.MaxDestinationStorageUsagePercent != nil {
		destinationUpdatedModel.MaxDestinationStorageUsagePercent = updateDto.MaxDestinationStorageUsagePercent
	}

	if updateDto.MaxConcurrentConnections != nil {
		destinationUpdatedModel.MaxConcurrentConnections = updateDto.MaxConcurrentConnections
	}

	if updateDto.TasksCount != nil {
		destinationUpdatedModel.TasksCount = *updateDto.TasksCount
	}

	if updateDto.TotalSpaceUsageBytes != nil {
		destinationUpdatedModel.TotalSpaceUsageBytes = updateDto.TotalSpaceUsageBytes.Uint64()
	}

	if updateDto.DownloadBytesSecRateLimit != nil {
		downloadBytesSecRateLimitUint := updateDto.DownloadBytesSecRateLimit.Uint64()
		destinationUpdatedModel.DownloadBytesSecRateLimit = &downloadBytesSecRateLimitUint
	}

	if updateDto.UploadBytesSecRateLimit != nil {
		uploadBytesSecRateLimitUint := updateDto.UploadBytesSecRateLimit.Uint64()
		destinationUpdatedModel.UploadBytesSecRateLimit = &uploadBytesSecRateLimitUint
	}

	if updateDto.ObjectStorageProvider != nil {
		objectStorageProviderStr := updateDto.ObjectStorageProvider.String()
		destinationUpdatedModel.ObjectStorageProvider = &objectStorageProviderStr
	}

	if updateDto.ObjectStorageProviderRegion != nil {
		objectStorageProviderRegionStr := updateDto.ObjectStorageProviderRegion.String()
		destinationUpdatedModel.ObjectStorageProviderRegion = &objectStorageProviderRegionStr
	}

	if updateDto.ObjectStorageProviderAccessKeyId != nil {
		objectStorageProviderAccessKeyIdStr := updateDto.ObjectStorageProviderAccessKeyId.String()
		destinationUpdatedModel.ObjectStorageProviderAccessKeyId = &objectStorageProviderAccessKeyIdStr
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
		destinationUpdatedModel.ObjectStorageProviderSecretAccessKey = &encryptedProviderSecretAccessKey
	}

	if updateDto.ObjectStorageEndpointUrl != nil {
		objectStorageEndpointUrlStr := updateDto.ObjectStorageEndpointUrl.String()
		destinationUpdatedModel.ObjectStorageEndpointUrl = &objectStorageEndpointUrlStr
	}

	if updateDto.ObjectStorageBucketName != nil {
		objectStorageBucketNameStr := updateDto.ObjectStorageBucketName.String()
		destinationUpdatedModel.ObjectStorageBucketName = &objectStorageBucketNameStr
	}

	if updateDto.RemoteHostType != nil {
		remoteHostTypeStr := updateDto.RemoteHostType.String()
		destinationUpdatedModel.RemoteHostType = &remoteHostTypeStr
	}

	if updateDto.RemoteHostname != nil {
		remoteHostnameStr := updateDto.RemoteHostname.String()
		destinationUpdatedModel.RemoteHostname = &remoteHostnameStr
	}

	if updateDto.RemoteHostNetworkPort != nil {
		remoteHostNetworkPortUint := updateDto.RemoteHostNetworkPort.Uint16()
		destinationUpdatedModel.RemoteHostNetworkPort = &remoteHostNetworkPortUint
	}

	if updateDto.RemoteHostUsername != nil {
		remoteHostUsernameStr := updateDto.RemoteHostUsername.String()
		destinationUpdatedModel.RemoteHostUsername = &remoteHostUsernameStr
	}

	if updateDto.RemoteHostPassword != nil {
		encryptedPassword, err := infraHelper.EncryptStr(
			dbEncryptSecret, updateDto.RemoteHostPassword.String(),
		)
		if err != nil {
			return errors.New("EncryptPasswordFailed: " + err.Error())
		}
		destinationUpdatedModel.RemoteHostPassword = &encryptedPassword
	}

	if updateDto.RemoteHostPrivateKeyFilePath != nil {
		remoteHostPrivateKeyFilePathStr := updateDto.RemoteHostPrivateKeyFilePath.String()
		destinationUpdatedModel.RemoteHostPrivateKeyFilePath = &remoteHostPrivateKeyFilePathStr
	}

	if updateDto.RemoteHostConnectionTimeoutSecs != nil {
		remoteHostConnectionTimeoutSecsUint := updateDto.RemoteHostConnectionTimeoutSecs.Uint64()
		destinationUpdatedModel.RemoteHostConnectionTimeoutSecs = &remoteHostConnectionTimeoutSecsUint
	}

	if updateDto.RemoteHostConnectionRetrySecs != nil {
		remoteHostConnectionRetrySecsUint := updateDto.RemoteHostConnectionRetrySecs.Uint64()
		destinationUpdatedModel.RemoteHostConnectionRetrySecs = &remoteHostConnectionRetrySecsUint
	}

	err := repo.persistentDbSvc.Handler.
		Model(&dbModel.BackupDestination{}).
		Where("id = ?", updateDto.DestinationId.Uint64()).
		Updates(&destinationUpdatedModel).Error
	if err != nil {
		return err
	}

	boolUpdateMap := map[string]interface{}{}
	if updateDto.SkipCertificateVerification != nil {
		boolUpdateMap["skip_certificate_verification"] = *updateDto.SkipCertificateVerification
	}

	if len(boolUpdateMap) == 0 {
		return nil
	}

	return repo.persistentDbSvc.Handler.
		Model(&dbModel.BackupDestination{}).
		Where("id = ?", updateDto.DestinationId.Uint64()).
		Updates(boolUpdateMap).Error
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

	timeoutSecs := useCase.BackupJobDefaultTimeoutSecs
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
		timeoutSecs.Uint64(), createDto.MaxTaskRetentionCount, createDto.MaxTaskRetentionDays,
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
		jobUpdatedModel.TimeoutSecs = updateDto.TimeoutSecs.Uint64()
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

	if updateDto.TasksCount != nil {
		jobUpdatedModel.TasksCount = *updateDto.TasksCount
	}

	if updateDto.TotalSpaceUsageBytes != nil {
		jobUpdatedModel.TotalSpaceUsageBytes = updateDto.TotalSpaceUsageBytes.Uint64()
	}

	if updateDto.LastRunAt != nil {
		lastRunAtTime := updateDto.LastRunAt.GetAsGoTime()
		jobUpdatedModel.LastRunAt = &lastRunAtTime
	}

	if updateDto.LastRunStatus != nil {
		lastRunStatusStr := updateDto.LastRunStatus.String()
		jobUpdatedModel.LastRunStatus = &lastRunStatusStr
	}

	if updateDto.NextRunAt != nil {
		nextRunAtTime := updateDto.NextRunAt.GetAsGoTime()
		jobUpdatedModel.NextRunAt = &nextRunAtTime
	}

	err := repo.persistentDbSvc.Handler.
		Model(&dbModel.BackupJob{}).
		Where("id = ?", updateDto.JobId.Uint64()).
		Updates(&jobUpdatedModel).Error
	if err != nil {
		return err
	}

	boolUpdateMap := map[string]interface{}{}
	if updateDto.JobStatus != nil {
		boolUpdateMap["job_status"] = *updateDto.JobStatus
	}

	if len(boolUpdateMap) > 0 {
		err = repo.persistentDbSvc.Handler.
			Model(&dbModel.BackupJob{}).
			Where("id = ?", updateDto.JobId.Uint64()).
			Updates(boolUpdateMap).Error
		if err != nil {
			return err
		}
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
	ContainerAccountIds    []string
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
			TimeoutSecs:       jobEntity.TimeoutSecs.Uint64(),
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
			ContainerAccountIds:    []string{},
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
	accountCmdRepo *infra.AccountCmdRepo,
	containerWithMetrics dto.ContainerWithMetrics,
	accountEntity entity.Account,
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

func (repo *BackupCmdRepo) createContainerArchive(
	containerImageCmdRepo *infra.ContainerImageCmdRepo,
	containerWithMetrics dto.ContainerWithMetrics,
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
	imageIdErrorTag := "[imageId/" + snapshotImageId.String() + "] "
	if err != nil {
		return archiveFile, errors.New(
			imageIdErrorTag + "CreateArchiveFailed: " + err.Error(),
		)
	}

	deleteSnapshotDto := dto.DeleteContainerImage{
		AccountId: containerWithMetrics.AccountId,
		ImageId:   snapshotImageId,
	}
	err = containerImageCmdRepo.Delete(deleteSnapshotDto)
	if err != nil {
		if !strings.Contains(err.Error(), "image is in use") {
			return archiveFile, errors.New(
				imageIdErrorTag + "DeleteSnapshotImageFailed: " + err.Error(),
			)
		}
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
	taskRunDetails.ExecutionOutput += "[containerId/" + containerIdStr + "] " + errorMessage + "\n"
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
	iDestinationEntity entity.IBackupDestination,
) (string, error) {
	unencryptedDestEnvPrefix := "RCLONE_CONFIG_RAWDEST"
	encryptedDestEnvPrefix := "RCLONE_CONFIG_ENCDEST"

	var backupBinaryEnvs []string

	var encryptionKey valueObject.Password
	switch destinationEntity := iDestinationEntity.(type) {
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
	switch destinationEntity := iDestinationEntity.(type) {
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

	backupBinaryFlags := []string{}
	maxConcurrentConnections := uint16(2)
	var uploadBytesSecRateLimitPtr *valueObject.Byte
	switch destinationEntity := taskRunDetails.DestinationEntity.(type) {
	case entity.BackupDestinationRemoteHost:
		if destinationEntity.MaxConcurrentConnections != nil {
			maxConcurrentConnections = *destinationEntity.MaxConcurrentConnections
		}

		if destinationEntity.UploadBytesSecRateLimit != nil {
			uploadBytesSecRateLimitPtr = destinationEntity.UploadBytesSecRateLimit
		}
	case entity.BackupDestinationObjectStorage:
		backupBinaryFlags = append(backupBinaryFlags, "--s3-no-check-bucket")
		if destinationEntity.MaxConcurrentConnections != nil {
			maxConcurrentConnections = *destinationEntity.MaxConcurrentConnections
		}

		if destinationEntity.UploadBytesSecRateLimit != nil {
			uploadBytesSecRateLimitPtr = destinationEntity.UploadBytesSecRateLimit
		}
	}
	backupBinaryFlags = append(
		backupBinaryFlags,
		"--transfers="+strconv.Itoa(int(maxConcurrentConnections)),
	)
	if uploadBytesSecRateLimitPtr != nil {
		backupBinaryFlags = append(
			backupBinaryFlags,
			"--bwlimit="+uploadBytesSecRateLimitPtr.String()+"B",
		)
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
	backupBinaryCmd := backupBinaryCli + " " + strings.Join(backupBinaryFlags, " ")
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

	accountIdStr := containerWithMetrics.AccountId.String()
	if !slices.Contains(taskRunDetails.ContainerAccountIds, accountIdStr) {
		taskRunDetails.ContainerAccountIds = append(
			taskRunDetails.ContainerAccountIds, containerWithMetrics.AccountId.String(),
		)
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

type BackupAccountContainersResponse struct {
	ExecOutputStr      string
	FailedContainerIds []valueObject.ContainerId
}

func (repo *BackupCmdRepo) backupAccountContainers(
	accountCmdRepo *infra.AccountCmdRepo,
	containerImageCmdRepo *infra.ContainerImageCmdRepo,
	preTaskAccountEntity entity.Account,
	containerWithMetricsSlice []dto.ContainerWithMetrics,
	jobTmpDir valueObject.UnixFilePath,
	taskIdRunDetailsMap map[valueObject.BackupTaskId]BackupTaskRunDetails,
) (backupResponse BackupAccountContainersResponse) {
	failedExecPrefix := "[accountId/" + preTaskAccountEntity.Id.String() + "] "

	localProcExecOutput := ""
	localProcFailedContainerIds := []valueObject.ContainerId{}

	for _, containerWithMetrics := range containerWithMetricsSlice {
		containerIdStr := containerWithMetrics.Id.String()
		failedContainerExecPrefix := failedExecPrefix + "[containerId/" + containerIdStr + "] "

		err := repo.accountStorageAllocator(
			accountCmdRepo, containerWithMetrics, preTaskAccountEntity,
		)
		if err != nil {
			localProcExecOutput += failedContainerExecPrefix +
				"AccountStorageAllocatorFailed: " + err.Error() + "\n"
			localProcFailedContainerIds = append(
				localProcFailedContainerIds, containerWithMetrics.Id,
			)
			continue
		}

		archiveEntity, err := repo.createContainerArchive(
			containerImageCmdRepo, containerWithMetrics, jobTmpDir,
		)
		if err != nil {
			localProcExecOutput += failedContainerExecPrefix +
				"CreateContainerArchiveFailed: " + err.Error() + "\n"
			localProcFailedContainerIds = append(
				localProcFailedContainerIds, containerWithMetrics.Id,
			)
			continue
		}

		for taskId, preUploadTaskRunDetails := range taskIdRunDetailsMap {
			externalProcTaskRunDetails := repo.uploadContainerArchive(
				preUploadTaskRunDetails, containerWithMetrics, archiveEntity,
			)
			taskIdRunDetailsMap[taskId] = externalProcTaskRunDetails
		}

		err = containerImageCmdRepo.DeleteArchive(archiveEntity)
		if err != nil {
			slog.Debug(
				"DeleteContainerArchiveFailed",
				slog.String("containerId", containerWithMetrics.Id.String()),
				slog.String("error", err.Error()),
			)
			continue
		}
	}

	return BackupAccountContainersResponse{
		ExecOutputStr:      localProcExecOutput,
		FailedContainerIds: localProcFailedContainerIds,
	}
}

func (repo *BackupCmdRepo) backupDestinationFreeSpaceValidator(
	jobEntity entity.BackupJob,
) (destinationIdsWithEnoughFreeSpace []valueObject.BackupDestinationId) {
	defaultMinFreeSpacePercent := uint8(5)

	userDataDirectoryStats, err := repo.readUserDataStats()
	if err != nil {
		slog.Debug(
			"ReadUserDataStatsFailed",
			slog.String("method", "FreeSpaceValidator"),
			slog.String("error", err.Error()),
		)
		return destinationIdsWithEnoughFreeSpace
	}

	for _, destinationId := range jobEntity.DestinationIds {
		iDestinationEntity, err := repo.backupQueryRepo.ReadFirstDestination(
			dto.ReadBackupDestinationsRequest{DestinationId: &destinationId}, false,
		)
		if err != nil {
			slog.Debug(
				"ReadBackupDestinationFailed",
				slog.String("method", "FreeSpaceValidator"),
				slog.String("error", err.Error()),
				slog.String("destinationId", destinationId.String()),
			)
			continue
		}

		desiredMinLocalFreeSpacePercent := defaultMinFreeSpacePercent
		var minLocalFreeSpacePercentPtr *uint8
		var maxDestinationStorageUsagePercentPtr *uint8
		var destinationPath valueObject.UnixFilePath

		switch destinationEntity := iDestinationEntity.(type) {
		case entity.BackupDestinationLocal:
			minLocalFreeSpacePercentPtr = destinationEntity.MinLocalStorageFreePercent
			maxDestinationStorageUsagePercentPtr = destinationEntity.MaxDestinationStorageUsagePercent
			destinationPath = destinationEntity.DestinationPath
		case entity.BackupDestinationRemoteHost:
			minLocalFreeSpacePercentPtr = destinationEntity.MinLocalStorageFreePercent
		case entity.BackupDestinationObjectStorage:
			minLocalFreeSpacePercentPtr = destinationEntity.MinLocalStorageFreePercent
		}
		if minLocalFreeSpacePercentPtr != nil {
			desiredMinLocalFreeSpacePercent = *minLocalFreeSpacePercentPtr
		}

		maxLocalSpaceUsagePercent := float64(100 - desiredMinLocalFreeSpacePercent)
		if userDataDirectoryStats.UsedPercent >= maxLocalSpaceUsagePercent {
			slog.Debug(
				"UserDataDirectorySpaceUsageExceeded",
				slog.Uint64("usedPercent", uint64(userDataDirectoryStats.UsedPercent)),
				slog.String("destinationId", destinationId.String()),
			)
			continue
		}

		if maxDestinationStorageUsagePercentPtr == nil {
			destinationIdsWithEnoughFreeSpace = append(destinationIdsWithEnoughFreeSpace, destinationId)
			continue
		}

		localDestinationDirectoryDiskStats, err := disk.Usage(destinationPath.String())
		if err != nil || localDestinationDirectoryDiskStats == nil {
			slog.Debug(
				"ReadDestinationDirectoryPathDiskStatsFailed",
				slog.String("destinationPath", destinationPath.String()),
				slog.String("destinationId", destinationId.String()),
			)
			continue
		}

		if uint8(localDestinationDirectoryDiskStats.UsedPercent) >= *maxDestinationStorageUsagePercentPtr {
			slog.Error(
				"DestinationDirectorySpaceUsageExceeded",
				slog.Uint64("usedPercent", uint64(localDestinationDirectoryDiskStats.UsedPercent)),
				slog.String("destinationId", destinationId.String()),
			)
			continue
		}

		destinationIdsWithEnoughFreeSpace = append(destinationIdsWithEnoughFreeSpace, destinationId)
	}

	return destinationIdsWithEnoughFreeSpace
}

func (repo *BackupCmdRepo) RunJob(
	runDto dto.RunBackupJob,
) (taskIds []valueObject.BackupTaskId, err error) {
	jobStartedAt := time.Now()

	jobEntity, err := repo.backupQueryRepo.ReadFirstJob(dto.ReadBackupJobsRequest{
		AccountId: &runDto.AccountId,
		JobId:     &runDto.JobId,
	})
	if err != nil {
		return taskIds, errors.New("ReadBackupJobFailed: " + err.Error())
	}

	destinationIdsWithEnoughFreeSpace := repo.backupDestinationFreeSpaceValidator(jobEntity)
	if len(destinationIdsWithEnoughFreeSpace) == 0 {
		return taskIds, errors.New("NoBackupDestinationsWithEnoughFreeSpace")
	}
	jobEntity.DestinationIds = destinationIdsWithEnoughFreeSpace

	accountIdContainerWithMetricsMap, err := repo.readAccountsContainersWithMetrics(jobEntity)
	if err != nil {
		return taskIds, errors.New("ReadAccountContainersFailed: " + err.Error())
	}

	taskIdRunDetailsMap := repo.backupTaskRunDetailsFactory(
		jobEntity, runDto.OperatorAccountId,
	)
	if len(taskIdRunDetailsMap) == 0 {
		return taskIds, errors.New("NoBackupTasksCreated")
	}

	rawJobTmpDir := fmt.Sprintf(
		"%s/backup/jobs/%d/%d",
		infraEnvs.NobodyDataDirectory, runDto.AccountId.Uint64(), runDto.JobId.Uint64(),
	)
	jobTmpDir, err := valueObject.NewUnixFilePath(rawJobTmpDir)
	if err != nil {
		return taskIds, errors.New("ValidateBackupJobTmpDirFailed: " + err.Error())
	}

	err = repo.createJobTmpDir(jobTmpDir)
	if err != nil {
		return taskIds, errors.New("CreateBackupJobTmpDirFailed: " + err.Error())
	}

	accountQueryRepo := infra.NewAccountQueryRepo(repo.persistentDbSvc)
	accountCmdRepo := infra.NewAccountCmdRepo(repo.persistentDbSvc)
	containerImageCmdRepo := infra.NewContainerImageCmdRepo(repo.persistentDbSvc)

	var localProcExecOutput strings.Builder
	localProcFailedContainerIds := []valueObject.ContainerId{}

	remainingJobRunTime := time.Duration(jobEntity.TimeoutSecs.Int64())
	accContainersBackupResponseChannel := make(chan BackupAccountContainersResponse)

	for accountId, containerWithMetricsSlice := range accountIdContainerWithMetricsMap {
		var localProcExecOutputWriteStr = func(message string) {
			localProcExecOutput.WriteString(
				"[accountId/" + accountId.String() + "] " + message + "\n",
			)
		}

		preTaskAccountEntity, err := accountQueryRepo.ReadById(accountId)
		if err != nil {
			localProcExecOutputWriteStr("ReadPreTaskAccountEntityFailed: " + err.Error())
			continue
		}

		accContainersBackupStartedAt := time.Now()
		go func() {
			accContainersBackupResponse := repo.backupAccountContainers(
				accountCmdRepo, containerImageCmdRepo, preTaskAccountEntity,
				containerWithMetricsSlice, jobTmpDir, taskIdRunDetailsMap,
			)
			accContainersBackupResponseChannel <- accContainersBackupResponse
		}()

		select {
		case accContainersBackupResponse := <-accContainersBackupResponseChannel:
			localProcExecOutput.WriteString(accContainersBackupResponse.ExecOutputStr)
			localProcFailedContainerIds = append(
				localProcFailedContainerIds,
				accContainersBackupResponse.FailedContainerIds...,
			)
		case <-time.After(remainingJobRunTime * time.Second):
			localProcExecOutputWriteStr("BackupTaskTimeout")
		}

		remainingJobRunTime -= time.Since(accContainersBackupStartedAt)
		if remainingJobRunTime <= 0 {
			break
		}

		postTaskAccountEntity, err := accountQueryRepo.ReadById(accountId)
		if err != nil {
			localProcExecOutputWriteStr("ReadPostTaskAccountEntityFailed: " + err.Error())
			continue
		}

		if preTaskAccountEntity.Quota.StorageBytes != postTaskAccountEntity.Quota.StorageBytes {
			err = accountCmdRepo.UpdateQuota(accountId, preTaskAccountEntity.Quota)
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

	localProcFailedContainerIdsStr := []string{}
	for _, containerId := range localProcFailedContainerIds {
		localProcFailedContainerIdsStr = append(localProcFailedContainerIdsStr, containerId.String())
	}

	for taskId, externalProcTaskRunDetails := range taskIdRunDetailsMap {
		taskIds = append(taskIds, taskId)

		combinedExecutionOutput := externalProcTaskRunDetails.ExecutionOutput
		localProcExecOutputStr := localProcExecOutput.String()
		if len(localProcExecOutputStr) > 0 {
			combinedExecutionOutput = localProcExecOutputStr + "\n" +
				externalProcTaskRunDetails.ExecutionOutput
		}

		combinedFailedContainerIds := append(
			externalProcTaskRunDetails.FailedContainerIds, localProcFailedContainerIdsStr...,
		)

		taskStatus := valueObject.BackupTaskStatusCompleted
		if len(combinedFailedContainerIds) > 0 {
			taskStatus = valueObject.BackupTaskStatusPartial
		}
		if len(externalProcTaskRunDetails.SuccessfulContainerIds) == 0 {
			taskStatus = valueObject.BackupTaskStatusFailed
		}

		if externalProcTaskRunDetails.StartedAt == externalProcTaskRunDetails.FinishedAt {
			externalProcTaskRunDetails.FinishedAt = time.Now()
			externalProcTaskRunDetails.ElapsedSecs = uint64(time.Since(jobStartedAt).Seconds())
		}

		taskModelUpdated := dbModel.BackupTask{
			TaskStatus:             taskStatus.String(),
			ExecutionOutput:        &combinedExecutionOutput,
			ContainerAccountIds:    externalProcTaskRunDetails.ContainerAccountIds,
			SuccessfulContainerIds: externalProcTaskRunDetails.SuccessfulContainerIds,
			FailedContainerIds:     combinedFailedContainerIds,
			SizeBytes:              &externalProcTaskRunDetails.SizeBytes,
			ElapsedSecs:            &externalProcTaskRunDetails.ElapsedSecs,
			FinishedAt:             &externalProcTaskRunDetails.FinishedAt,
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

	return taskIds, nil
}

func (repo *BackupCmdRepo) toggleNobodyUserSession(sessionStatus bool) error {
	_, _ = infraHelper.RunCmd("pkill", "-u", "nobody")
	_, _ = infraHelper.RunCmd("pkill", "-u", "nobody")
	time.Sleep(3 * time.Second)

	userModFlags := []string{"--add-subuids", "--add-subgids"}
	if !sessionStatus {
		userModFlags = []string{"--del-subuids", "--del-subgids"}
	}

	uidGidsRange := "100000-165535"
	_, err := infraHelper.RunCmd(
		"usermod", userModFlags[0], uidGidsRange, userModFlags[1], uidGidsRange, "nobody",
	)
	if err != nil {
		return errors.New("UsermodAddSubUidsNobodyFailed: " + err.Error())
	}

	if sessionStatus {
		return infraHelper.EnableLingering(valueObject.NobodyAccountId)
	}

	return infraHelper.DisableLingering(valueObject.NobodyAccountId)
}

func (repo *BackupCmdRepo) restoreImageArchive(
	imageArchiveEntity entity.ContainerImageArchive,
	containerImageQueryRepo *infra.ContainerImageQueryRepo,
	containerImageCmdRepo *infra.ContainerImageCmdRepo,
) (imageEntity entity.ContainerImage, err error) {
	if imageArchiveEntity.AccountId == valueObject.SystemAccountId {
		imageArchiveEntity.AccountId = valueObject.NobodyAccountId
	}

	importImageArchiveDto := dto.ImportContainerImageArchive{
		AccountId:       imageArchiveEntity.AccountId,
		ArchiveFilePath: &imageArchiveEntity.UnixFilePath,
	}
	imageId, err := containerImageCmdRepo.ImportArchive(importImageArchiveDto)
	if err != nil {
		return imageEntity, errors.New("ImportContainerImageArchiveFailed: " + err.Error())
	}

	containerImageEntity, err := containerImageQueryRepo.ReadById(
		imageArchiveEntity.AccountId, imageId,
	)
	if err != nil {
		return imageEntity, errors.New("ReadContainerImageFailed: " + err.Error())
	}

	if containerImageEntity.OriginContainerDetails == nil {
		return imageEntity, errors.New("OriginContainerDetailsNotFound")
	}

	return containerImageEntity, nil
}

func (repo *BackupCmdRepo) restoreContainerArchive(
	archiveEntity entity.ContainerImageArchive,
	accountQueryRepo *infra.AccountQueryRepo,
	accountCmdRepo *infra.AccountCmdRepo,
	containerQueryRepo *infra.ContainerQueryRepo,
	containerCmdRepo *infra.ContainerCmdRepo,
	containerImageQueryRepo *infra.ContainerImageQueryRepo,
	containerImageCmdRepo *infra.ContainerImageCmdRepo,
	containerProfileQueryRepo *infra.ContainerProfileQueryRepo,
	containerProxyCmdRepo *infra.ContainerProxyCmdRepo,
	mappingQueryRepo *infra.MappingQueryRepo,
	mappingCmdRepo *infra.MappingCmdRepo,
	activityRecordCmdRepo *infra.ActivityRecordCmdRepo,
	shouldReplaceExistingContainers bool,
	shouldRestoreMappings bool,
) (containerId valueObject.ContainerId, err error) {
	containerImageEntity, err := repo.restoreImageArchive(
		archiveEntity, containerImageQueryRepo, containerImageCmdRepo,
	)
	if err != nil {
		return containerId, errors.New("InitialImageArchiveRestoreFailed: " + err.Error())
	}

	originalAccountOwnerId := containerImageEntity.OriginContainerDetails.AccountId
	if containerImageEntity.AccountId != originalAccountOwnerId {
		_, err = accountQueryRepo.ReadById(originalAccountOwnerId)
		if err == nil {
			deleteImageDto := dto.DeleteContainerImage{
				AccountId: containerImageEntity.AccountId,
				ImageId:   containerImageEntity.Id,
			}
			err = containerImageCmdRepo.Delete(deleteImageDto)
			if err != nil {
				return containerId, errors.New("DeleteInitialContainerImageFailed: " + err.Error())
			}

			archiveEntity.AccountId = originalAccountOwnerId
			containerImageEntity, err = repo.restoreImageArchive(
				archiveEntity, containerImageQueryRepo, containerImageCmdRepo,
			)
			if err != nil {
				return containerId, errors.New(
					"PostOwnershipFixImageArchiveRestoreFailed: " + err.Error(),
				)
			}
		}
	}

	containerEntity, err := containerQueryRepo.ReadFirst(dto.ReadContainersRequest{
		ContainerId: []valueObject.ContainerId{
			containerImageEntity.OriginContainerDetails.Id,
		},
	})
	previousContainerStillExists := err == nil
	if previousContainerStillExists && shouldReplaceExistingContainers {
		deleteContainerDto := dto.DeleteContainer{
			AccountId:   containerEntity.AccountId,
			ContainerId: containerEntity.Id,
		}

		err = containerCmdRepo.Delete(deleteContainerDto)
		if err != nil {
			return containerId, errors.New("DeleteExistingContainerFailed: " + err.Error())
		}
	}

	rawContainerHostname := containerImageEntity.OriginContainerDetails.Hostname.String()
	if previousContainerStillExists && !shouldReplaceExistingContainers {
		archiveCreatedAtStr := archiveEntity.CreatedAt.String()
		rawContainerHostname = archiveCreatedAtStr + ".restored." + rawContainerHostname
	}
	containerHostname, err := valueObject.NewFqdn(rawContainerHostname)
	if err != nil {
		return containerId, errors.New("ValidateContainerHostnameFailed: " + err.Error())
	}

	removableEnvsRegex := regexp.MustCompile(`^(PRIMARY_VHOST|HOSTNAME)`)
	adjustedEnvs := []valueObject.ContainerEnv{}
	for _, originalEnv := range containerImageEntity.OriginContainerDetails.Envs {
		if removableEnvsRegex.MatchString(originalEnv.String()) {
			continue
		}
		adjustedEnvs = append(adjustedEnvs, originalEnv)
	}

	createContainerDto := dto.CreateContainer{
		AccountId:          archiveEntity.AccountId,
		Hostname:           containerHostname,
		ImageAddress:       containerImageEntity.ImageAddress,
		ImageId:            &containerImageEntity.Id,
		PortBindings:       containerImageEntity.OriginContainerDetails.PortBindings,
		Envs:               adjustedEnvs,
		Entrypoint:         containerImageEntity.Entrypoint,
		ProfileId:          &containerImageEntity.OriginContainerDetails.ProfileId,
		RestartPolicy:      &containerImageEntity.OriginContainerDetails.RestartPolicy,
		AutoCreateMappings: shouldRestoreMappings,
	}

	containerId, err = useCase.CreateContainer(
		containerQueryRepo, containerCmdRepo, containerImageQueryRepo,
		containerImageCmdRepo, accountQueryRepo, accountCmdRepo, containerProfileQueryRepo,
		mappingQueryRepo, mappingCmdRepo, containerProxyCmdRepo, activityRecordCmdRepo,
		createContainerDto,
	)
	if err != nil {
		return containerId, errors.New("RestoreContainerFailed: " + err.Error())
	}

	err = os.Remove(archiveEntity.UnixFilePath.String())
	if err != nil {
		actualArchiveFilePath, err := containerImageCmdRepo.ImageArchiveFileLocator(
			archiveEntity.UnixFilePath,
		)
		if err != nil {
			return containerId, errors.New("LocateContainerImageArchiveFailed: " + err.Error())
		}
		err = os.Remove(actualArchiveFilePath.String())
		if err != nil {
			return containerId, errors.New("DeleteContainerImageArchiveFailed: " + err.Error())
		}
	}

	return containerId, nil
}

func (repo *BackupCmdRepo) RestoreTask(
	requestRestoreDto dto.RestoreBackupTaskRequest,
) (responseRestoreDto dto.RestoreBackupTaskResponse, err error) {
	responseRestoreDto.SuccessfulContainerIds = []valueObject.ContainerId{}
	responseRestoreDto.FailedContainerImageIds = []valueObject.ContainerImageId{}

	taskArchiveProvided := requestRestoreDto.ArchiveId != nil
	if !taskArchiveProvided {
		createArchiveDto := dto.CreateBackupTaskArchive{
			TaskId:                    *requestRestoreDto.TaskId,
			TimeoutSecs:               requestRestoreDto.TimeoutSecs,
			ContainerAccountIds:       requestRestoreDto.ContainerAccountIds,
			ContainerIds:              requestRestoreDto.ContainerIds,
			ExceptContainerAccountIds: requestRestoreDto.ExceptContainerAccountIds,
			ExceptContainerIds:        requestRestoreDto.ExceptContainerIds,
			OperatorAccountId:         requestRestoreDto.OperatorAccountId,
		}

		archiveId, err := repo.CreateTaskArchive(createArchiveDto)
		if err != nil {
			return responseRestoreDto, errors.New("CreateTaskArchiveFailed: " + err.Error())
		}
		requestRestoreDto.ArchiveId = &archiveId
	}

	archiveEntity, err := repo.backupQueryRepo.ReadFirstTaskArchive(
		dto.ReadBackupTaskArchivesRequest{ArchiveId: requestRestoreDto.ArchiveId},
	)
	if err != nil {
		return responseRestoreDto, errors.New("ReadBackupTaskArchiveFailed: " + err.Error())
	}

	userDataDirectoryStats, err := repo.readUserDataStats()
	if err != nil {
		return responseRestoreDto, errors.New("ReadUserDataDirStatsError: " + err.Error())
	}
	// @see https://ntorga.com/gzip-bzip2-xz-zstd-7z-brotli-or-lz4/
	necessaryFreeStorageBytes := archiveEntity.SizeBytes.Uint64() * 3
	if necessaryFreeStorageBytes > userDataDirectoryStats.Free {
		return responseRestoreDto, errors.New("InsufficientUserDataDirectoryFreeSpace")
	}

	restoreBaseTmpDir, err := valueObject.NewUnixFilePath(infraEnvs.RestoreBackupTaskTmpDir)
	if err != nil {
		return responseRestoreDto, errors.New("ValidateRestoreBaseTaskTmpDirFailed: " + err.Error())
	}
	restoreBaseTmpDirStr := restoreBaseTmpDir.String()
	err = infraHelper.MakeDir(restoreBaseTmpDirStr)
	if err != nil {
		return responseRestoreDto, errors.New("MakeRestoreBaseTaskTmpDirFailed: " + err.Error())
	}

	_, err = infraHelper.RunCmd(
		"tar", "-xf", archiveEntity.UnixFilePath.String(), "-C", restoreBaseTmpDirStr,
	)
	if err != nil {
		return responseRestoreDto, errors.New("ExtractTaskArchiveFailed: " + err.Error())
	}

	if !taskArchiveProvided {
		err = os.Remove(archiveEntity.UnixFilePath.String())
		if err != nil {
			return responseRestoreDto, errors.New("DeleteTempTaskArchiveFailed: " + err.Error())
		}
	}

	_, err = infraHelper.RunCmd("chown", "-R", "nobody:nogroup", restoreBaseTmpDirStr)
	if err != nil {
		return responseRestoreDto, errors.New("ChownRestoreBaseTaskTmpDirFailed: " + err.Error())
	}

	rawRestoreTaskTmpDir := restoreBaseTmpDirStr + "/" +
		archiveEntity.UnixFilePath.ReadFileNameWithoutExtension().String()
	restoreTaskTmpDir, err := valueObject.NewUnixFilePath(rawRestoreTaskTmpDir)
	if err != nil {
		return responseRestoreDto, errors.New("ValidateRestoreTmpDirFailed: " + err.Error())
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
		return responseRestoreDto, errors.New("ReadContainerImageArchivesFailed: " + err.Error())
	}

	if len(containerImagesResponseDto.Archives) == 0 {
		return responseRestoreDto, errors.New("NoContainerImageArchivesFound")
	}

	shouldReplaceExistingContainers := false
	if requestRestoreDto.ShouldReplaceExistingContainers != nil {
		shouldReplaceExistingContainers = *requestRestoreDto.ShouldReplaceExistingContainers
	}

	shouldRestoreMappings := true
	if requestRestoreDto.ShouldRestoreMappings != nil {
		shouldRestoreMappings = *requestRestoreDto.ShouldRestoreMappings
	}

	accountQueryRepo := infra.NewAccountQueryRepo(repo.persistentDbSvc)
	accountCmdRepo := infra.NewAccountCmdRepo(repo.persistentDbSvc)
	containerQueryRepo := infra.NewContainerQueryRepo(repo.persistentDbSvc)
	containerCmdRepo := infra.NewContainerCmdRepo(repo.persistentDbSvc)
	containerImageCmdRepo := infra.NewContainerImageCmdRepo(repo.persistentDbSvc)
	containerProfileQueryRepo := infra.NewContainerProfileQueryRepo(repo.persistentDbSvc)
	containerProxyCmdRepo := infra.NewContainerProxyCmdRepo(repo.persistentDbSvc)
	mappingQueryRepo := infra.NewMappingQueryRepo(repo.persistentDbSvc)
	mappingCmdRepo := infra.NewMappingCmdRepo(repo.persistentDbSvc)
	activityRecordCmdRepo := infra.NewActivityRecordCmdRepo(repo.trailDbSvc)

	err = repo.toggleNobodyUserSession(true)
	if err != nil {
		return responseRestoreDto, errors.New("ToggleNobodyUserSessionFailed: " + err.Error())
	}

	for _, imageArchiveEntity := range containerImagesResponseDto.Archives {
		containerId, err := repo.restoreContainerArchive(
			imageArchiveEntity, accountQueryRepo, accountCmdRepo, containerQueryRepo,
			containerCmdRepo, containerImageQueryRepo, containerImageCmdRepo,
			containerProfileQueryRepo, containerProxyCmdRepo, mappingQueryRepo, mappingCmdRepo,
			activityRecordCmdRepo, shouldReplaceExistingContainers, shouldRestoreMappings,
		)
		if err != nil {
			slog.Debug(
				"RestoreContainerArchiveFailed",
				slog.String("imageArchiveEntityId", imageArchiveEntity.ImageId.String()),
				slog.String("error", err.Error()),
			)
			responseRestoreDto.FailedContainerImageIds = append(
				responseRestoreDto.FailedContainerImageIds, imageArchiveEntity.ImageId,
			)
			continue
		}
		responseRestoreDto.SuccessfulContainerIds = append(
			responseRestoreDto.SuccessfulContainerIds, containerId,
		)
	}

	err = os.RemoveAll(restoreTaskTmpDir.String())
	if err != nil {
		return responseRestoreDto, errors.New("RemoveRestoreTmpDirFailed: " + err.Error())
	}

	err = repo.toggleNobodyUserSession(false)
	if err != nil {
		return responseRestoreDto, errors.New("ToggleNobodyUserSessionFailed: " + err.Error())
	}

	return responseRestoreDto, nil
}

func (repo *BackupCmdRepo) killJobPid(
	jobId valueObject.BackupJobId,
	elapsedSecs uint64,
) error {
	elapsedSecsSafetyMargin := uint64(60)
	if elapsedSecs > elapsedSecsSafetyMargin {
		elapsedSecs -= elapsedSecsSafetyMargin
	}
	elapsedSecsStr := strconv.Itoa(int(elapsedSecs))

	findJobPidsCmd := `ps -e -o "pid,etimes,cmd:128" --sort=etimes --no-headers |` +
		` awk '/backup job run/ && /--job-id ` + jobId.String() + `/ &&` +
		` $2 >= ` + elapsedSecsStr + ` {print $1}'`

	rawJobPidsOutput, err := infraHelper.RunCmdWithSubShell(findJobPidsCmd)
	if err != nil {
		return errors.New("FindBackupJobPidError: " + err.Error())
	}

	if len(rawJobPidsOutput) == 0 {
		return errors.New("BackupJobPidsNotFound")
	}

	rawJobPidsOutputSlice := strings.Split(rawJobPidsOutput, "\n")

	jobPidsStrSlice := []string{}
	for _, rawJobPidStr := range rawJobPidsOutputSlice {
		rawJobPidStr = strings.TrimSpace(rawJobPidStr)
		if len(rawJobPidStr) == 0 {
			continue
		}

		jobPidInt, err := strconv.Atoi(rawJobPidStr)
		if err != nil {
			continue
		}

		jobPidsStrSlice = append(jobPidsStrSlice, strconv.Itoa(jobPidInt))
	}

	if len(jobPidsStrSlice) == 0 {
		return errors.New("ValidBackupJobPidsNotFound")
	}

	_, err = infraHelper.RunCmdWithSubShell("kill -9 " + strings.Join(jobPidsStrSlice, " "))
	return err
}

func (repo *BackupCmdRepo) UpdateTask(updateDto dto.UpdateBackupTask) error {
	taskEntity, err := repo.backupQueryRepo.ReadFirstTask(
		dto.ReadBackupTasksRequest{TaskId: &updateDto.TaskId},
	)
	if err != nil {
		return errors.New("ReadBackupTaskFailed: " + err.Error())
	}

	if updateDto.TaskStatus == nil {
		return nil
	}

	elapsedSecs := uint64(0)
	if taskEntity.StartedAt != nil {
		startedAt := taskEntity.StartedAt.GetAsGoTime()
		elapsedSecs = uint64(time.Since(startedAt).Seconds())
	}

	if *updateDto.TaskStatus == valueObject.BackupTaskStatusCancelled {
		err = repo.killJobPid(taskEntity.JobId, elapsedSecs)
		if err != nil {
			slog.Debug("KillJobPidFailed", slog.String("error", err.Error()))
		}
	}

	finishedAt := time.Now()

	taskUpdatedModel := dbModel.BackupTask{
		ID:          taskEntity.TaskId.Uint64(),
		TaskStatus:  updateDto.TaskStatus.String(),
		ElapsedSecs: &elapsedSecs,
		FinishedAt:  &finishedAt,
	}

	return repo.persistentDbSvc.Handler.Model(&dbModel.BackupTask{}).
		Where("id = ?", taskEntity.TaskId.Uint64()).
		Updates(&taskUpdatedModel).Error
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

	iDestinationEntity, err := repo.backupQueryRepo.ReadFirstDestination(
		dto.ReadBackupDestinationsRequest{DestinationId: &taskEntity.DestinationId}, true,
	)
	if err != nil {
		return archiveId, errors.New("ReadBackupDestinationFailed: " + err.Error())
	}

	necessaryFreeStorageBytes = 0
	containerArchiveEntities, err := repo.readRemoteContainerArchives(
		taskEntity, iDestinationEntity, createDto.ContainerAccountIds, createDto.ContainerIds,
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

	accountHomeDirStr := infraEnvs.NobodyDataDirectory
	if createDto.OperatorAccountId != valueObject.SystemAccountId {
		accountQueryRepo := infra.NewAccountQueryRepo(repo.persistentDbSvc)
		operatorAccountEntity, err := accountQueryRepo.ReadById(createDto.OperatorAccountId)
		if err != nil {
			return archiveId, errors.New("ReadOperatorAccountFailed: " + err.Error())
		}
		if necessaryFreeStorageBytes*2 >= operatorAccountEntity.Quota.StorageBytes.Uint64() {
			return archiveId, errors.New("InsufficientOperatorAccountQuota")
		}

		accountHomeDirStr = operatorAccountEntity.HomeDirectory.String()
	}
	archivesDirectoryStr := accountHomeDirStr + TasksArchivesRelativePath

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

	backupBinaryFlags := []string{}
	maxConcurrentConnections := uint16(2)
	var downloadBytesSecRateLimitPtr *valueObject.Byte
	switch destinationEntity := iDestinationEntity.(type) {
	case entity.BackupDestinationRemoteHost:
		if destinationEntity.MaxConcurrentConnections != nil {
			maxConcurrentConnections = *destinationEntity.MaxConcurrentConnections
		}

		if destinationEntity.DownloadBytesSecRateLimit != nil {
			downloadBytesSecRateLimitPtr = destinationEntity.DownloadBytesSecRateLimit
		}
	case entity.BackupDestinationObjectStorage:
		backupBinaryFlags = append(backupBinaryFlags, "--s3-no-check-bucket")
		if destinationEntity.MaxConcurrentConnections != nil {
			maxConcurrentConnections = *destinationEntity.MaxConcurrentConnections
		}

		if destinationEntity.DownloadBytesSecRateLimit != nil {
			downloadBytesSecRateLimitPtr = destinationEntity.DownloadBytesSecRateLimit
		}
	}
	backupBinaryFlags = append(
		backupBinaryFlags,
		"--transfers="+strconv.Itoa(int(maxConcurrentConnections)),
	)
	if downloadBytesSecRateLimitPtr != nil {
		backupBinaryFlags = append(
			backupBinaryFlags,
			"--bwlimit="+downloadBytesSecRateLimitPtr.String()+"B",
		)
	}

	backupBinaryCli, err := repo.backupBinaryCliFactory(iDestinationEntity)
	if err != nil {
		return archiveId, errors.New("BackupCliFactoryFailed: " + err.Error())
	}
	backupBinaryCli += " " + strings.Join(backupBinaryFlags, " ")

	taskRemotePathStr := repo.readTaskRemotePath(iDestinationEntity, taskEntity.TaskId)

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

	return os.Remove(taskArchiveEntity.UnixFilePath.String())
}
