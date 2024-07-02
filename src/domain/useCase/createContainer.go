package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
)

func CreateContainer(
	containerQueryRepo repository.ContainerQueryRepo,
	containerCmdRepo repository.ContainerCmdRepo,
	accountQueryRepo repository.AccountQueryRepo,
	accountCmdRepo repository.AccountCmdRepo,
	containerProfileQueryRepo repository.ContainerProfileQueryRepo,
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	containerProxyCmdRepo repository.ContainerProxyCmdRepo,
	createDto dto.CreateContainer,
) error {
	err := CheckAccountQuota(
		accountQueryRepo,
		createDto.AccountId,
		containerProfileQueryRepo,
		*createDto.ProfileId,
		nil,
	)
	if err != nil {
		log.Printf("QuotaCheckError: %s", err)
		return err
	}

	_, err = containerQueryRepo.ReadByHostname(createDto.Hostname)
	if err == nil {
		log.Printf("ContainerHostnameAlreadyExists: %s", createDto.Hostname)
		return errors.New("ContainerHostnameAlreadyExists")
	}

	containerId, err := containerCmdRepo.Create(createDto)
	if err != nil {
		log.Printf("CreateContainerError: %s", err)
		return errors.New("CreateContainerInfraError")
	}

	err = accountCmdRepo.UpdateQuotaUsage(createDto.AccountId)
	if err != nil {
		log.Printf("UpdateAccountQuotaError: %s", err)
		return errors.New("UpdateAccountQuotaError")
	}

	log.Printf(
		"ContainerId '%s' (%s) created for AccountId '%s'.",
		containerId.String(),
		createDto.ImageAddress.String(),
		createDto.AccountId.String(),
	)

	if createDto.ImageAddress.IsSpeediaOs() {
		err = containerProxyCmdRepo.Create(containerId)
		if err != nil {
			log.Printf("CreateContainerProxyError: %s", err)
			return errors.New("CreateContainerProxyInfraError")
		}
	}

	if !createDto.AutoCreateMappings {
		return nil
	}

	err = CreateMappingsWithContainerId(
		containerQueryRepo,
		mappingQueryRepo,
		mappingCmdRepo,
		containerProxyCmdRepo,
		containerId,
	)
	if err != nil {
		log.Printf("CreateAutoMappingsError: %s", err)
	}

	return nil
}
