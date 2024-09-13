package useCase

import (
	"errors"
	"log/slog"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
)

func DeleteContainerImage(
	containerImageQueryRepo repository.ContainerImageQueryRepo,
	containerImageCmdRepo repository.ContainerImageCmdRepo,
	containerQueryRepo repository.ContainerQueryRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteContainerImage,
) error {
	_, err := containerImageQueryRepo.ReadById(deleteDto.AccountId, deleteDto.ImageId)
	if err != nil {
		return errors.New("ContainerImageNotFound")
	}

	containersUsingImage, err := containerQueryRepo.ReadByImageId(
		deleteDto.AccountId, deleteDto.ImageId,
	)
	if err != nil {
		slog.Error("ReadContainerByImageIdInfraError", slog.Any("error", err))
		return errors.New("ReadContainerByImageIdInfraError")
	}
	if len(containersUsingImage) > 0 {
		return errors.New("ContainerImageInUseCannotBeDeleted")
	}

	err = containerImageCmdRepo.Delete(deleteDto)
	if err != nil {
		slog.Error("DeleteContainerImageInfraError", slog.Any("error", err))
		return errors.New("DeleteContainerImageError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		DeleteContainerImage(deleteDto)

	return nil
}
