package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
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

	readContainersDto := dto.ReadContainersRequest{
		Pagination:       ContainersDefaultPagination,
		ContainerImageId: &deleteDto.ImageId,
	}

	responseDto, err := ReadContainers(containerQueryRepo, readContainersDto)
	if err != nil {
		return errors.New("ReadContainersInfraError")
	}

	if len(responseDto.Containers) > 0 {
		return errors.New("CannotDeleteContainerImageInUse")
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
