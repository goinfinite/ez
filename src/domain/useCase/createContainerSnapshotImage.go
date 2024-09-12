package useCase

import (
	"errors"
	"log/slog"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
)

func CreateContainerSnapshotImage(
	containerImageCmdRepo repository.ContainerImageCmdRepo,
	containerQueryRepo repository.ContainerQueryRepo,
	accountQueryRepo repository.AccountQueryRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateContainerSnapshotImage,
) error {
	containerEntityWithMetrics, err := containerQueryRepo.ReadWithMetricsById(
		createDto.AccountId, createDto.ContainerId,
	)
	if err != nil {
		return errors.New("ContainerNotFound")
	}

	accountEntity, err := accountQueryRepo.ReadById(createDto.AccountId)
	if err != nil {
		slog.Error("ReadAccountInfoInfraError", slog.Any("error", err))
		return errors.New("ReadAccountInfoInfraError")
	}

	storageAvailable := accountEntity.Quota.StorageBytes - accountEntity.QuotaUsage.StorageBytes
	isThereStorageAvailable := storageAvailable >= containerEntityWithMetrics.Metrics.StorageSpaceBytes
	if !isThereStorageAvailable {
		return errors.New("AccountStorageQuotaUsageExceeded")
	}

	imageId, err := containerImageCmdRepo.CreateSnapshot(createDto)
	if err != nil {
		slog.Error("CreateContainerSnapshotImageInfraError", slog.Any("error", err))
		return errors.New("CreateContainerSnapshotImageInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateContainerSnapshotImage(createDto, imageId)

	return nil
}
