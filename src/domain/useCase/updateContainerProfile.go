package useCase

import (
	"errors"
	"log/slog"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
)

func updateContainersAfterProfileUpdate(
	containerQueryRepo repository.ContainerQueryRepo,
	containerCmdRepo repository.ContainerCmdRepo,
	updateDto dto.UpdateContainerProfile,
) error {
	containers, err := containerQueryRepo.Read()
	if err != nil {
		slog.Error("ReadContainersInfraError", slog.Any("error", err))
		return errors.New("ReadContainersInfraError")
	}

	for _, container := range containers {
		if container.ProfileId != updateDto.ProfileId {
			continue
		}

		updateContainerDto := dto.NewUpdateContainer(
			container.AccountId, container.Id, &container.Status,
			&updateDto.ProfileId, updateDto.OperatorAccountId,
			updateDto.OperatorIpAddress,
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

func UpdateContainerProfile(
	containerProfileQueryRepo repository.ContainerProfileQueryRepo,
	containerProfileCmdRepo repository.ContainerProfileCmdRepo,
	containerQueryRepo repository.ContainerQueryRepo,
	containerCmdRepo repository.ContainerCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	updateDto dto.UpdateContainerProfile,
) error {
	_, err := containerProfileQueryRepo.ReadById(updateDto.ProfileId)
	if err != nil {
		return errors.New("ContainerProfileNotFound")
	}

	err = containerProfileCmdRepo.Update(updateDto)
	if err != nil {
		slog.Error("UpdateContainerProfileInfraError", slog.Any("error", err))
		return errors.New("UpdateContainerProfileInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).UpdateContainerProfile(updateDto)

	shouldUpdateContainers := updateDto.BaseSpecs != nil
	if !shouldUpdateContainers {
		return nil
	}

	err = updateContainersAfterProfileUpdate(
		containerQueryRepo, containerCmdRepo, updateDto,
	)
	if err != nil {
		slog.Error("UpdateContainersAfterProfileUpdate", slog.Any("error", err))
		return errors.New("UpdateContainersAfterProfileUpdate")
	}

	return nil
}
