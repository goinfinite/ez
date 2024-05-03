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
) (valueObject.AccessTokenValue, error) {
	var accessToken valueObject.AccessTokenValue

	containerEntity, err := containerQueryRepo.GetById(containerId)
	if err != nil {
		log.Printf("ContainerNotFound: %s", err)
		return accessToken, errors.New("ContainerNotFound")
	}

	if containerEntity.IsSpeediaOs() {
		log.Printf("ContainerIsNotSpeediaOs: %s", containerEntity.ImageAddress)
		return accessToken, errors.New("ContainerIsNotSpeediaOs")
	}

	if !containerEntity.IsRunning() {
		return accessToken, errors.New("ContainerIsNotRunning")
	}

	accessToken, err = containerCmdRepo.GenerateContainerSessionToken(containerId)
	if err != nil {
		log.Printf("GenerateContainerSessionTokenError: %s", err)
		return accessToken, errors.New("GenerateContainerSessionTokenError")
	}

	return accessToken, nil
}
