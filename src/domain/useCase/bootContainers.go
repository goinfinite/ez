package useCase

import (
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

func BootContainers(
	containerQueryRepo repository.ContainerQueryRepo,
	containerCmdRepo repository.ContainerCmdRepo,
) {
	paginationDto := dto.Pagination{
		PageNumber:   0,
		ItemsPerPage: 1000,
	}

	restartPolicy, _ := valueObject.NewContainerRestartPolicy("always")
	readContainersDto := dto.ReadContainersRequest{
		Pagination:             paginationDto,
		ContainerRestartPolicy: &restartPolicy,
	}

	responseDto, err := ReadContainers(containerQueryRepo, readContainersDto)
	if err != nil {
		slog.Debug("ReadContainersDuringBootError", slog.Any("error", err))
		return
	}

	for _, currentContainer := range responseDto.Containers {
		newContainerStatus := true
		updateDto := dto.NewUpdateContainer(
			currentContainer.AccountId, currentContainer.Id, &newContainerStatus,
			nil, valueObject.SystemAccountId, valueObject.SystemIpAddress,
		)

		err = containerCmdRepo.Update(updateDto)
		if err != nil {
			slog.Debug(
				"StartContainerDuringBootInfraError",
				slog.String("containerId", currentContainer.Id.String()),
				slog.Any("error", err),
			)
			continue
		}
	}
}
