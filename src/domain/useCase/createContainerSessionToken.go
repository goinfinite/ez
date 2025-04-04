package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

func CreateContainerSessionToken(
	containerQueryRepo repository.ContainerQueryRepo,
	containerCmdRepo repository.ContainerCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateContainerSessionToken,
) (accessToken valueObject.AccessTokenValue, err error) {
	readContainersDto := dto.ReadContainersRequest{
		Pagination:  ContainersDefaultPagination,
		ContainerId: []valueObject.ContainerId{createDto.ContainerId},
	}

	responseDto, err := ReadContainers(containerQueryRepo, readContainersDto)
	if err != nil || len(responseDto.Containers) == 0 {
		return accessToken, errors.New("ContainerNotFound")
	}
	containerEntity := responseDto.Containers[0]

	if !containerEntity.ImageAddress.IsInfiniteOs() {
		return accessToken, errors.New("NotInfiniteOs")
	}

	if !containerEntity.IsRunning() {
		return accessToken, errors.New("ContainerIsNotRunning")
	}

	accessToken, err = containerCmdRepo.CreateContainerSessionToken(createDto)
	if err != nil {
		slog.Error("CreateContainerSessionTokenInfraError", slog.Any("error", err))
		return accessToken, errors.New("CreateContainerSessionTokenInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateContainerSessionToken(createDto)

	return accessToken, nil
}
