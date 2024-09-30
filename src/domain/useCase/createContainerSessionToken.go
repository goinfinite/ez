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
	containerEntity, err := containerQueryRepo.ReadById(createDto.ContainerId)
	if err != nil {
		return accessToken, errors.New("ContainerNotFound")
	}

	if !containerEntity.ImageAddress.IsSpeediaOs() {
		return accessToken, errors.New("NotSpeediaOs")
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
