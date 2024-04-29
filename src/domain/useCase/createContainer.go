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
	accQueryRepo repository.AccQueryRepo,
	accCmdRepo repository.AccCmdRepo,
	containerProfileQueryRepo repository.ContainerProfileQueryRepo,
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	createDto dto.CreateContainer,
) error {
	err := CheckAccountQuota(
		accQueryRepo,
		createDto.AccountId,
		containerProfileQueryRepo,
		*createDto.ProfileId,
		nil,
	)
	if err != nil {
		log.Printf("QuotaCheckError: %s", err)
		return err
	}

	_, err = containerQueryRepo.GetByHostname(createDto.Hostname)
	if err == nil {
		log.Printf("ContainerHostnameAlreadyExists: %s", createDto.Hostname)
		return errors.New("ContainerHostnameAlreadyExists")
	}

	containerId, err := containerCmdRepo.Create(createDto)
	if err != nil {
		log.Printf("CreateContainerError: %s", err)
		return errors.New("CreateContainerInfraError")
	}

	err = accCmdRepo.UpdateQuotaUsage(createDto.AccountId)
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

	if !createDto.AutoCreateMappings {
		return nil
	}

	return AddMappingsWithContainerId(
		containerQueryRepo,
		mappingQueryRepo,
		mappingCmdRepo,
		containerId,
	)
}