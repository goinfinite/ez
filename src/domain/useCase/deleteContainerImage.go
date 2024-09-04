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
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteContainerImage,
) error {
	_, err := containerImageQueryRepo.ReadById(deleteDto.AccountId, deleteDto.ImageId)
	if err != nil {
		return errors.New("ContainerImageNotFound")
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
