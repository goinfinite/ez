package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/repository"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

func CreateContainerSnapshotImage(
	containerImageQueryRepo repository.ContainerImageQueryRepo,
	containerImageCmdRepo repository.ContainerImageCmdRepo,
	containerQueryRepo repository.ContainerQueryRepo,
	accountQueryRepo repository.AccountQueryRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createSnapshotDto dto.CreateContainerSnapshotImage,
) (imageEntity entity.ContainerImage, err error) {
	containerEntityWithMetrics, err := containerQueryRepo.ReadFirstWithMetrics(
		dto.ReadContainersRequest{
			ContainerId: []valueObject.ContainerId{createSnapshotDto.ContainerId},
		},
	)
	if err != nil {
		return imageEntity, errors.New("ContainerNotFound")
	}

	containerAccountId := containerEntityWithMetrics.AccountId
	accountEntity, err := accountQueryRepo.ReadFirst(dto.ReadAccountsRequest{
		AccountId: &containerAccountId,
	})
	if err != nil {
		slog.Error("ReadAccountInfoError", slog.String("error", err.Error()))
		return imageEntity, errors.New("ReadAccountInfoInfraError")
	}

	storageAvailable := accountEntity.Quota.StorageBytes - accountEntity.QuotaUsage.StorageBytes
	isThereStorageAvailable := storageAvailable >= containerEntityWithMetrics.Metrics.StorageSpaceBytes
	if !isThereStorageAvailable {
		return imageEntity, errors.New("AccountStorageQuotaUsageExceeded")
	}

	imageId, err := containerImageCmdRepo.CreateSnapshot(createSnapshotDto)
	if err != nil {
		slog.Error("CreateContainerSnapshotImageError", slog.String("error", err.Error()))
		return imageEntity, errors.New("CreateContainerSnapshotImageInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateContainerSnapshotImage(createSnapshotDto, containerAccountId, imageId)

	imageEntity, err = containerImageQueryRepo.ReadById(containerAccountId, imageId)
	if err != nil {
		slog.Error("ReadContainerSnapshotImageError", slog.String("error", err.Error()))
		return imageEntity, errors.New("ReadContainerSnapshotImageInfraError")
	}

	if createSnapshotDto.ShouldCreateArchive == nil {
		return imageEntity, nil
	}

	if !*createSnapshotDto.ShouldCreateArchive {
		return imageEntity, nil
	}

	createArchiveDto := dto.NewCreateContainerImageArchive(
		containerAccountId, imageId, createSnapshotDto.ArchiveCompressionFormat,
		createSnapshotDto.ArchiveDestinationPath, createSnapshotDto.OperatorAccountId,
		createSnapshotDto.OperatorIpAddress,
	)
	_, err = CreateContainerImageArchive(
		containerImageQueryRepo, containerImageCmdRepo,
		accountQueryRepo, activityRecordCmdRepo, createArchiveDto,
	)
	if err != nil {
		return imageEntity, err
	}

	if createSnapshotDto.ShouldDiscardImage == nil {
		return imageEntity, nil
	}

	if !*createSnapshotDto.ShouldDiscardImage {
		return imageEntity, nil
	}

	deleteImageDto := dto.NewDeleteContainerImage(
		containerAccountId, imageId,
		createSnapshotDto.OperatorAccountId, createSnapshotDto.OperatorIpAddress,
	)
	err = containerImageCmdRepo.Delete(deleteImageDto)
	if err != nil {
		slog.Error("DeleteContainerSnapshotImageInfraError", slog.String("error", err.Error()))
		return imageEntity, errors.New("DeleteContainerSnapshotImageInfraError")
	}

	return imageEntity, nil
}
