package dbModel

import (
	"errors"
	"time"

	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type BackupDestination struct {
	ID                                   uint64 `gorm:"primarykey"`
	AccountID                            uint64 `gorm:"not null"`
	Name                                 string `gorm:"not null"`
	Description                          *string
	Type                                 string `gorm:"not null"`
	Path                                 string `gorm:"not null"`
	MinLocalStorageFreePercent           *uint8
	MaxDestinationStorageUsagePercent    *uint8
	TotalSpaceUsageBytes                 *uint64
	TotalSpaceUsagePercent               *uint8
	MaxConcurrentConnections             *uint16
	DownloadBytesSecRateLimit            *uint64
	UploadBytesSecRateLimit              *uint64
	SkipCertificateVerification          *bool
	ObjectStorageProvider                *string
	ObjectStorageProviderRegion          *string
	ObjectStorageProviderAccessKeyId     *string
	ObjectStorageProviderSecretAccessKey *string
	ObjectStorageEndpointUrl             *string
	ObjectStorageBucketName              *string
	RemoteHostType                       *string
	RemoteHostname                       *string
	RemoteHostNetworkPort                *uint16
	RemoteHostUsername                   *string
	RemoteHostPassword                   *string
	RemoteHostPrivateKeyFilePath         *string
	RemoteHostConnectionTimeoutSecs      *uint16
	RemoteHostConnectionRetrySecs        *uint16
	CreatedAt                            time.Time `gorm:"not null"`
	UpdatedAt                            time.Time `gorm:"not null"`
}

func (model BackupDestination) TableName() string {
	return "backup_destinations"
}

func NewBackupDestination(
	id, accountId uint64,
	name string,
	description *string,
	destinationType, path string,
	minLocalStorageFreePercent, maxDestinationStorageUsagePercent *uint8,
	maxConcurrentConnections *uint16,
	downloadBytesSecRateLimit, uploadBytesSecRateLimit *uint64,
	skipCertificateVerification *bool,
	objectStorageProvider, objectStorageProviderRegion, objectStorageProviderAccessKeyId,
	objectStorageProviderSecretAccessKey, objectStorageEndpointUrl, objectStorageBucketName,
	remoteHostType, remoteHostname, remoteHostUsername, remoteHostPassword, remoteHostPrivateKeyFilePath *string,
	remoteHostNetworkPort, remoteHostConnectionTimeoutSecs, remoteHostConnectionRetrySecs *uint16,
) BackupDestination {
	destinationModel := BackupDestination{
		ID:                                   id,
		AccountID:                            accountId,
		Name:                                 name,
		Description:                          description,
		Type:                                 destinationType,
		Path:                                 path,
		MinLocalStorageFreePercent:           minLocalStorageFreePercent,
		MaxDestinationStorageUsagePercent:    maxDestinationStorageUsagePercent,
		MaxConcurrentConnections:             maxConcurrentConnections,
		DownloadBytesSecRateLimit:            downloadBytesSecRateLimit,
		UploadBytesSecRateLimit:              uploadBytesSecRateLimit,
		SkipCertificateVerification:          skipCertificateVerification,
		ObjectStorageProvider:                objectStorageProvider,
		ObjectStorageProviderRegion:          objectStorageProviderRegion,
		ObjectStorageProviderAccessKeyId:     objectStorageProviderAccessKeyId,
		ObjectStorageProviderSecretAccessKey: objectStorageProviderSecretAccessKey,
		ObjectStorageEndpointUrl:             objectStorageEndpointUrl,
		ObjectStorageBucketName:              objectStorageBucketName,
		RemoteHostType:                       remoteHostType,
		RemoteHostname:                       remoteHostname,
		RemoteHostNetworkPort:                remoteHostNetworkPort,
		RemoteHostUsername:                   remoteHostUsername,
		RemoteHostPassword:                   remoteHostPassword,
		RemoteHostPrivateKeyFilePath:         remoteHostPrivateKeyFilePath,
		RemoteHostConnectionTimeoutSecs:      remoteHostConnectionTimeoutSecs,
		RemoteHostConnectionRetrySecs:        remoteHostConnectionRetrySecs,
	}

	if id != 0 {
		destinationModel.ID = id
	}

	return destinationModel
}

func (model BackupDestination) ToEntity() (
	destinationEntity entity.IBackupDestination, err error,
) {
	destinationId, err := valueObject.NewBackupDestinationId(model.ID)
	if err != nil {
		return destinationEntity, err
	}

	accountId, err := valueObject.NewAccountId(model.AccountID)
	if err != nil {
		return destinationEntity, err
	}

	destinationName, err := valueObject.NewBackupDestinationName(model.Name)
	if err != nil {
		return destinationEntity, err
	}

	var destinationDescriptionPtr *valueObject.BackupDestinationDescription
	if model.Description != nil {
		destinationDescription, err := valueObject.NewBackupDestinationDescription(*model.Description)
		if err != nil {
			return destinationEntity, err
		}
		destinationDescriptionPtr = &destinationDescription
	}

	destinationPath, err := valueObject.NewUnixFilePath(model.Path)
	if err != nil {
		return destinationEntity, err
	}

	var totalSpaceUsageBytesPtr *valueObject.Byte
	if model.TotalSpaceUsageBytes != nil {
		totalSpaceUsageBytes, err := valueObject.NewByte(*model.TotalSpaceUsageBytes)
		if err != nil {
			return destinationEntity, err
		}
		totalSpaceUsageBytesPtr = &totalSpaceUsageBytes
	}

	destinationType, err := valueObject.NewBackupDestinationType(model.Type)
	if err != nil {
		return destinationEntity, err
	}

	backupDestinationBase := entity.BackupDestinationBase{
		DestinationId:          destinationId,
		AccountId:              accountId,
		DestinationName:        destinationName,
		DestinationDescription: destinationDescriptionPtr,
		DestinationType:        destinationType,
		DestinationPath:        destinationPath,
		TotalSpaceUsageBytes:   totalSpaceUsageBytesPtr,
		TotalSpaceUsagePercent: model.TotalSpaceUsagePercent,
		CreatedAt:              valueObject.NewUnixTimeWithGoTime(model.CreatedAt),
		UpdatedAt:              valueObject.NewUnixTimeWithGoTime(model.UpdatedAt),
	}

	backupDestinationRemoteBase := entity.BackupDestinationRemoteBase{
		BackupDestinationBase:       backupDestinationBase,
		MaxConcurrentConnections:    model.MaxConcurrentConnections,
		DownloadBytesSecRateLimit:   model.DownloadBytesSecRateLimit,
		UploadBytesSecRateLimit:     model.UploadBytesSecRateLimit,
		SkipCertificateVerification: model.SkipCertificateVerification,
	}

	switch destinationType {
	case valueObject.BackupDestinationTypeLocal:
		return entity.BackupDestinationLocal{
			BackupDestinationBase: backupDestinationBase,
		}, nil
	case valueObject.BackupDestinationTypeObjectStorage:
		objectStorageProvider, err := valueObject.NewObjectStorageProvider(
			*model.ObjectStorageProvider,
		)
		if err != nil {
			return destinationEntity, err
		}

		var objectStorageProviderRegionPtr *valueObject.ObjectStorageProviderRegion
		if model.ObjectStorageProviderRegion != nil {
			objectStorageProviderRegion, err := valueObject.NewObjectStorageProviderRegion(
				*model.ObjectStorageProviderRegion,
			)
			if err != nil {
				return destinationEntity, err
			}
			objectStorageProviderRegionPtr = &objectStorageProviderRegion
		}

		objectStorageProviderAccessKeyId, err := valueObject.NewObjectStorageProviderAccessKeyId(
			*model.ObjectStorageProviderAccessKeyId,
		)
		if err != nil {
			return destinationEntity, err
		}

		objectStorageEndpointUrl, err := valueObject.NewUrl(*model.ObjectStorageEndpointUrl)
		if err != nil {
			return destinationEntity, err
		}

		objectStorageBucketName, err := valueObject.NewObjectStorageBucketName(
			*model.ObjectStorageBucketName,
		)
		if err != nil {
			return destinationEntity, err
		}

		return entity.BackupDestinationObjectStorage{
			BackupDestinationRemoteBase:      backupDestinationRemoteBase,
			ObjectStorageProvider:            &objectStorageProvider,
			ObjectStorageProviderRegion:      objectStorageProviderRegionPtr,
			ObjectStorageProviderAccessKeyId: &objectStorageProviderAccessKeyId,
			ObjectStorageEndpointUrl:         &objectStorageEndpointUrl,
			ObjectStorageBucketName:          &objectStorageBucketName,
		}, nil
	case valueObject.BackupDestinationTypeRemoteHost:
		remoteHostType, err := valueObject.NewBackupDestinationRemoteHostType(
			*model.RemoteHostType,
		)
		if err != nil {
			return destinationEntity, err
		}

		remoteHostname, err := valueObject.NewNetworkHost(*model.RemoteHostname)
		if err != nil {
			return destinationEntity, err
		}

		remoteHostNetworkPort, err := valueObject.NewNetworkPort(*model.RemoteHostNetworkPort)
		if err != nil {
			return destinationEntity, err
		}

		remoteHostUsername, err := valueObject.NewUnixUsername(*model.RemoteHostUsername)
		if err != nil {
			return destinationEntity, err
		}

		var remoteHostPrivateKeyFilePathPtr *valueObject.UnixFilePath
		if model.RemoteHostPrivateKeyFilePath != nil {
			remoteHostPrivateKeyFilePath, err := valueObject.NewUnixFilePath(
				*model.RemoteHostPrivateKeyFilePath,
			)
			if err != nil {
				return destinationEntity, err
			}
			remoteHostPrivateKeyFilePathPtr = &remoteHostPrivateKeyFilePath
		}

		return entity.BackupDestinationRemoteHost{
			BackupDestinationRemoteBase:     backupDestinationRemoteBase,
			RemoteHostType:                  &remoteHostType,
			RemoteHostname:                  &remoteHostname,
			RemoteHostNetworkPort:           &remoteHostNetworkPort,
			RemoteHostUsername:              &remoteHostUsername,
			RemoteHostPrivateKeyFilePath:    remoteHostPrivateKeyFilePathPtr,
			RemoteHostConnectionTimeoutSecs: model.RemoteHostConnectionTimeoutSecs,
			RemoteHostConnectionRetrySecs:   model.RemoteHostConnectionRetrySecs,
		}, nil
	default:
		return destinationEntity, errors.New("UnsupportedBackupDestinationType")
	}
}
