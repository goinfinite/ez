package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
)

func CreateContainerSnapshotImage(
	containerImageQueryRepo repository.ContainerImageQueryRepo,
	containerImageCmdRepo repository.ContainerImageCmdRepo,
	containerQueryRepo repository.ContainerQueryRepo,
	accountQueryRepo repository.AccountQueryRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createSnapshotDto dto.CreateContainerSnapshotImage,
) error {
	withMetrics := true
	readContainersDto := dto.ReadContainersRequest{
		Pagination:  ContainersDefaultPagination,
		ContainerId: &createSnapshotDto.ContainerId,
		WithMetrics: &withMetrics,
	}

	responseDto, err := ReadContainers(containerQueryRepo, readContainersDto)
	if err != nil || len(responseDto.ContainersWithMetrics) == 0 {
		return errors.New("ContainerNotFound")
	}
	containerEntityWithMetrics := responseDto.ContainersWithMetrics[0]
	containerAccountId := containerEntityWithMetrics.AccountId

	accountEntity, err := accountQueryRepo.ReadById(containerAccountId)
	if err != nil {
		slog.Error("ReadAccountInfoInfraError", slog.Any("error", err))
		return errors.New("ReadAccountInfoInfraError")
	}

	storageAvailable := accountEntity.Quota.StorageBytes - accountEntity.QuotaUsage.StorageBytes
	isThereStorageAvailable := storageAvailable >= containerEntityWithMetrics.Metrics.StorageSpaceBytes
	if !isThereStorageAvailable {
		return errors.New("AccountStorageQuotaUsageExceeded")
	}

	imageId, err := containerImageCmdRepo.CreateSnapshot(createSnapshotDto)
	if err != nil {
		slog.Error("CreateContainerSnapshotImageInfraError", slog.Any("error", err))
		return errors.New("CreateContainerSnapshotImageInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateContainerSnapshotImage(createSnapshotDto, containerAccountId, imageId)

	if createSnapshotDto.ShouldCreateArchive == nil {
		return nil
	}

	if !*createSnapshotDto.ShouldCreateArchive {
		return nil
	}

	createArchiveDto := dto.NewCreateContainerImageArchiveFile(
		containerAccountId, imageId, createSnapshotDto.ArchiveCompressionFormat, nil,
		createSnapshotDto.OperatorAccountId, createSnapshotDto.OperatorIpAddress,
	)
	_, err = CreateContainerImageArchiveFile(
		containerImageQueryRepo, containerImageCmdRepo,
		accountQueryRepo, activityRecordCmdRepo, createArchiveDto,
	)
	if err != nil {
		return err
	}

	if createSnapshotDto.ShouldDiscardImage == nil {
		return nil
	}

	if !*createSnapshotDto.ShouldDiscardImage {
		return nil
	}

	deleteImageDto := dto.NewDeleteContainerImage(
		containerAccountId, imageId,
		createSnapshotDto.OperatorAccountId, createSnapshotDto.OperatorIpAddress,
	)
	err = containerImageCmdRepo.Delete(deleteImageDto)
	if err != nil {
		slog.Error("DeleteContainerSnapshotImageInfraError", slog.Any("error", err))
		return errors.New("DeleteContainerSnapshotImageInfraError")
	}

	return nil
}
