package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

func CreateBackupDestination(
	backupCmdRepo repository.BackupCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateBackupDestination,
) error {
	switch createDto.DestinationType {
	case valueObject.BackupDestinationTypeObjectStorage:
		if createDto.ObjectStorageProvider == nil {
			return errors.New("ObjectStorageProviderIsRequired")
		}
		if createDto.ObjectStorageProviderAccessKeyId == nil {
			return errors.New("ObjectStorageProviderAccessKeyIdIsRequired")
		}
		if createDto.ObjectStorageProviderSecretAccessKey == nil {
			return errors.New("ObjectStorageProviderSecretAccessKeyIsRequired")
		}
		if createDto.ObjectStorageBucketName == nil {
			return errors.New("ObjectStorageBucketNameIsRequired")
		}
		if *createDto.ObjectStorageProvider == valueObject.ObjectStorageProviderCustom &&
			createDto.ObjectStorageEndpointUrl == nil {
			return errors.New("ObjectStorageEndpointUrlIsRequired")
		}
	case valueObject.BackupDestinationTypeRemoteHost:
		if createDto.RemoteHostType == nil {
			return errors.New("RemoteHostTypeIsRequired")
		}
		if createDto.RemoteHostname == nil {
			return errors.New("RemoteHostnameIsRequired")
		}
		if createDto.RemoteHostNetworkPort == nil {
			return errors.New("RemoteHostNetworkPortIsRequired")
		}
		if createDto.RemoteHostUsername == nil {
			return errors.New("RemoteHostUsernameIsRequired")
		}
		if createDto.RemoteHostPassword == nil && createDto.RemoteHostPrivateKeyFilePath == nil {
			return errors.New("PasswordOrPrivateKeyFilePathIsRequired")
		}
	}

	backupDestinationId, err := backupCmdRepo.CreateDestination(createDto)
	if err != nil {
		slog.Error("CreateBackupDestinationInfraError", slog.Any("error", err))
		return errors.New("CreateBackupDestinationInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateBackupDestination(createDto, backupDestinationId)

	return nil
}
