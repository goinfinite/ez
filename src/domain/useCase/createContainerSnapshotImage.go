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
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateContainerSnapshotImage,
) error {
	_, err := containerQueryRepo.ReadById(createDto.ContainerId)
	if err != nil {
		return errors.New("ContainerNotFound")
	}

	imageId, err := containerImageCmdRepo.CreateSnapshot(createDto)
	if err != nil {
		slog.Error(
			"CreateContainerSnapshotImageInfraError", slog.Any("error", err),
		)
		return errors.New("CreateContainerSnapshotImageInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateContainerSnapshotImage(createDto, imageId)
	return nil
}
