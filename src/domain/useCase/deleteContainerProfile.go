package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/repository"
)

func resetContainersProfile(
	containerQueryRepo repository.ContainerQueryRepo,
	containerCmdRepo repository.ContainerCmdRepo,
	deleteDto dto.DeleteContainerProfile,
) error {
	readContainersDto := dto.ReadContainersRequest{
		Pagination:         ContainersDefaultPagination,
		ContainerProfileId: &deleteDto.ProfileId,
	}

	responseDto, err := ReadContainers(containerQueryRepo, readContainersDto)
	if err != nil {
		return errors.New("ReadContainersInfraError")
	}

	defaultContainerProfile := entity.DefaultContainerProfile()
	for _, container := range responseDto.Containers {
		updateContainerDto := dto.NewUpdateContainer(
			container.AccountId, container.Id, &container.Status,
			&defaultContainerProfile.Id, deleteDto.OperatorAccountId,
			deleteDto.OperatorIpAddress,
		)

		err := containerCmdRepo.Update(updateContainerDto)
		if err != nil {
			slog.Debug(
				"UpdateContainerInfraError",
				slog.String("containerId", container.Id.String()),
				slog.Any("error", err),
			)
			continue
		}
	}

	return nil
}

func DeleteContainerProfile(
	containerProfileQueryRepo repository.ContainerProfileQueryRepo,
	containerProfileCmdRepo repository.ContainerProfileCmdRepo,
	containerQueryRepo repository.ContainerQueryRepo,
	containerCmdRepo repository.ContainerCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteContainerProfile,
) error {
	_, err := containerProfileQueryRepo.ReadById(deleteDto.ProfileId)
	if err != nil {
		return errors.New("ContainerProfileNotFound")
	}

	defaultContainerProfile := entity.DefaultContainerProfile()
	if deleteDto.ProfileId == defaultContainerProfile.Id {
		return errors.New("CannotDeleteDefaultContainerProfile")
	}

	err = resetContainersProfile(containerQueryRepo, containerCmdRepo, deleteDto)
	if err != nil {
		slog.Error("ResetContainersProfileInfraError", slog.Any("error", err))
		return errors.New("ResetContainersProfileInfraError")
	}

	err = containerProfileCmdRepo.Delete(deleteDto)
	if err != nil {
		slog.Error("DeleteContainerProfileInfraError", slog.Any("error", err))
		return errors.New("DeleteContainerProfileInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		DeleteContainerProfile(deleteDto)

	return nil
}
