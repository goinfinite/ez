package backupInfra

import (
	"errors"
	"os"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/valueObject"
	"github.com/goinfinite/ez/src/infra/db"
	dbModel "github.com/goinfinite/ez/src/infra/db/model"
	infraHelper "github.com/goinfinite/ez/src/infra/helper"
)

type BackupCmdRepo struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewBackupCmdRepo(
	persistentDbSvc *db.PersistentDatabaseService,
) *BackupCmdRepo {
	return &BackupCmdRepo{persistentDbSvc: persistentDbSvc}
}

func (repo *BackupCmdRepo) CreateDestination(
	createDto dto.CreateBackupDestination,
) (destinationId valueObject.BackupDestinationId, err error) {
	var descriptionPtr *string
	if createDto.DestinationDescription != nil {
		description := createDto.DestinationDescription.String()
		descriptionPtr = &description
	}

	destinationPath := "/"
	if createDto.DestinationPath != nil {
		destinationPath = createDto.DestinationPath.String()
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

	encryptSecretKey := os.Getenv("BACKUP_KEYS_SECRET")
	if encryptSecretKey == "" {
		return destinationId, errors.New("BackupKeysSecretMissing")
	}

	if createDto.ObjectStorageProviderSecretAccessKey != nil {
		encryptedProviderSecretAccessKey, err := infraHelper.EncryptStr(
			encryptSecretKey, createDto.ObjectStorageProviderSecretAccessKey.String(),
		)
		if err != nil {
			return destinationId, errors.New("EncryptProviderSecretAccessKeyFailed: " + err.Error())
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
			encryptSecretKey, createDto.RemoteHostPassword.String(),
		)
		if err != nil {
			return destinationId, errors.New("EncryptPasswordFailed: " + err.Error())
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
		descriptionPtr, createDto.DestinationType.String(), destinationPath,
		createDto.MinLocalStorageFreePercent, createDto.MaxDestinationStorageUsagePercent,
		createDto.MaxConcurrentConnections, createDto.DownloadBytesSecRateLimit,
		createDto.UploadBytesSecRateLimit, createDto.SkipCertificateVerification,
		objectStorageProviderPtr, objectStorageProviderRegionPtr, objectStorageProviderAccessKeyIdPtr,
		objectStorageProviderSecretAccessKeyPtr, objectStorageEndpointUrlPtr,
		objectStorageBucketNamePtr, remoteHostTypePtr, remoteHostnamePtr,
		remoteHostUsernamePtr, remoteHostPasswordPtr, remoteHostPrivateKeyFilePathPtr,
		remoteHostNetworkPortPtr, createDto.RemoteHostConnectionTimeoutSecs,
		createDto.RemoteHostConnectionRetrySecs,
	)

	err = repo.persistentDbSvc.Handler.Create(&destinationModel).Error
	if err != nil {
		return destinationId, err
	}

	return valueObject.NewBackupDestinationId(destinationModel.ID)
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

	encryptSecretKey := os.Getenv("BACKUP_KEYS_SECRET")
	if encryptSecretKey == "" {
		return errors.New("BackupKeysSecretMissing")
	}

	if updateDto.ObjectStorageProviderSecretAccessKey != nil {
		encryptedProviderSecretAccessKey, err := infraHelper.EncryptStr(
			encryptSecretKey, updateDto.ObjectStorageProviderSecretAccessKey.String(),
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
			encryptSecretKey, updateDto.RemoteHostPassword.String(),
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

	var ignoreContainerAccountIdsUint64 []uint64
	for _, ignoreContainerAccountId := range createDto.IgnoreContainerAccountIds {
		ignoreContainerAccountIdsUint64 = append(ignoreContainerAccountIdsUint64, ignoreContainerAccountId.Uint64())
	}

	var ignoreContainerIds []string
	for _, ignoreContainerId := range createDto.IgnoreContainerIds {
		ignoreContainerIds = append(ignoreContainerIds, ignoreContainerId.String())
	}

	jobModel := dbModel.NewBackupJob(
		0, createDto.AccountId.Uint64(), true, jobDescriptionPtr, destinationIdsUint64,
		retentionStrategy.String(), createDto.BackupSchedule.String(), archiveCompressionFormat.String(),
		timeoutSecs, createDto.MaxTaskRetentionCount, createDto.MaxTaskRetentionDays,
		createDto.MaxConcurrentCpuCores, containerAccountIdsUint64, containerIds,
		ignoreContainerAccountIdsUint64, ignoreContainerIds,
	)

	err = repo.persistentDbSvc.Handler.Create(&jobModel).Error
	if err != nil {
		return backupJobId, err
	}

	return valueObject.NewBackupJobId(jobModel.ID)
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

	if updateDto.IgnoreContainerAccountIds != nil {
		ignoreContainerAccountIdsUint64 := []uint64{}
		for _, ignoreContainerAccountId := range updateDto.IgnoreContainerAccountIds {
			ignoreContainerAccountIdsUint64 = append(ignoreContainerAccountIdsUint64, ignoreContainerAccountId.Uint64())
		}
		jobUpdatedModel.IgnoreContainerAccountIds = ignoreContainerAccountIdsUint64
	}

	if updateDto.IgnoreContainerIds != nil {
		ignoreContainerIds := []string{}
		for _, ignoreContainerId := range updateDto.IgnoreContainerIds {
			ignoreContainerIds = append(ignoreContainerIds, ignoreContainerId.String())
		}
		jobUpdatedModel.IgnoreContainerIds = ignoreContainerIds
	}

	return repo.persistentDbSvc.Handler.
		Model(&dbModel.BackupJob{}).
		Where("id = ?", updateDto.JobId.Uint64()).
		Updates(&jobUpdatedModel).Error
}

func (repo *BackupCmdRepo) DeleteJob(
	deleteDto dto.DeleteBackupJob,
) error {
	return repo.persistentDbSvc.Handler.Model(&dbModel.BackupJob{}).Delete(
		"id = ? AND account_id = ?",
		deleteDto.JobId.Uint64(), deleteDto.AccountId.Uint64(),
	).Error
}

func (repo *BackupCmdRepo) RunJob(runDto dto.RunBackupJob) error {
	return nil
}
