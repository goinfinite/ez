package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/repository"
)

func CreateContainerImageArchive(
	containerImageQueryRepo repository.ContainerImageQueryRepo,
	containerImageCmdRepo repository.ContainerImageCmdRepo,
	accountQueryRepo repository.AccountQueryRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateContainerImageArchive,
) (archiveFile entity.ContainerImageArchive, err error) {
	imageEntity, err := containerImageQueryRepo.ReadById(
		createDto.AccountId, createDto.ImageId,
	)
	if err != nil {
		return archiveFile, errors.New("ContainerImageNotFound")
	}

	accountEntity, err := accountQueryRepo.ReadById(createDto.AccountId)
	if err != nil {
		slog.Error("ReadAccountInfoInfraError", slog.Any("error", err))
		return archiveFile, errors.New("ReadAccountInfoInfraError")
	}

	accountStorageAvailable := accountEntity.Quota.StorageBytes - accountEntity.QuotaUsage.StorageBytes
	isThereStorageAvailable := accountStorageAvailable >= imageEntity.SizeBytes
	if !isThereStorageAvailable {
		return archiveFile, errors.New("AccountStorageQuotaUsageExceeded")
	}

	archiveFile, err = containerImageCmdRepo.CreateArchive(createDto)
	if err != nil {
		slog.Error("CreateContainerImageArchiveInfraError", slog.Any("error", err))
		return archiveFile, errors.New("CreateContainerImageArchiveInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateContainerImageArchive(createDto)

	return archiveFile, nil
}
