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
	createDto dto.CreateBackupDestinationRequest,
) (responseDto dto.CreateBackupDestinationResponse, err error) {
	switch createDto.DestinationType {
	case valueObject.BackupDestinationTypeObjectStorage:
		if createDto.ObjectStorageProvider == nil {
			return responseDto, errors.New("ObjectStorageProviderIsRequired")
		}
		if createDto.ObjectStorageProviderAccessKeyId == nil {
			return responseDto, errors.New("ObjectStorageProviderAccessKeyIdIsRequired")
		}
		if createDto.ObjectStorageProviderSecretAccessKey == nil {
			return responseDto, errors.New("ObjectStorageProviderSecretAccessKeyIsRequired")
		}
		if createDto.ObjectStorageBucketName == nil {
			return responseDto, errors.New("ObjectStorageBucketNameIsRequired")
		}
		if *createDto.ObjectStorageProvider == valueObject.ObjectStorageProviderCustom &&
			createDto.ObjectStorageEndpointUrl == nil {
			return responseDto, errors.New("ObjectStorageEndpointUrlIsRequired")
		}
	case valueObject.BackupDestinationTypeRemoteHost:
		if createDto.RemoteHostType == nil {
			return responseDto, errors.New("RemoteHostTypeIsRequired")
		}
		if createDto.RemoteHostname == nil {
			return responseDto, errors.New("RemoteHostnameIsRequired")
		}
		if createDto.RemoteHostNetworkPort == nil {
			return responseDto, errors.New("RemoteHostNetworkPortIsRequired")
		}
		if createDto.RemoteHostUsername == nil {
			return responseDto, errors.New("RemoteHostUsernameIsRequired")
		}
		if createDto.RemoteHostPassword == nil && createDto.RemoteHostPrivateKeyFilePath == nil {
			return responseDto, errors.New("PasswordOrPrivateKeyFilePathIsRequired")
		}
	}

	responseDto, err = backupCmdRepo.CreateDestination(createDto)
	if err != nil {
		slog.Error("CreateBackupDestinationInfraError", slog.Any("error", err))
		return responseDto, errors.New("CreateBackupDestinationInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateBackupDestination(createDto, responseDto.DestinationId)

	return responseDto, nil
}
