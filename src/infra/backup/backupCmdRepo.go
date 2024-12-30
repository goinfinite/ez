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

	encryptSecretKey := os.Getenv("BACKUP_KEYS_SECRET")
	if encryptSecretKey == "" {
		return destinationId, errors.New("BackupKeysSecretMissing")
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

	destinationId, err = valueObject.NewBackupDestinationId(destinationModel.ID)
	if err != nil {
		return destinationId, err
	}

	return destinationId, nil
}
