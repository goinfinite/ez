package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func ContainerAutoLogin(
	containerQueryRepo repository.ContainerQueryRepo,
	containerCmdRepo repository.ContainerCmdRepo,
	containerId valueObject.ContainerId,
) (accessToken valueObject.AccessTokenValue, err error) {
	containerEntity, err := containerQueryRepo.GetById(containerId)
	if err != nil {
		return accessToken, errors.New("ContainerNotFound")
	}

	if !containerEntity.ImageAddress.IsSpeediaOs() {
		return accessToken, errors.New("NotSpeediaOs")
	}

	if !containerEntity.IsRunning() {
		return accessToken, errors.New("ContainerIsNotRunning")
	}

	accessToken, err = containerCmdRepo.GenerateContainerSessionToken(containerId)
	if err != nil {
		log.Printf("GenerateContainerSessionTokenError: %s", err)
		return accessToken, errors.New("GenerateContainerSessionTokenInfraError")
	}

	return accessToken, nil
}
