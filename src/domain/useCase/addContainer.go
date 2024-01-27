package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
)

func AddContainer(
	containerQueryRepo repository.ContainerQueryRepo,
	containerCmdRepo repository.ContainerCmdRepo,
	accQueryRepo repository.AccQueryRepo,
	accCmdRepo repository.AccCmdRepo,
	containerProfileQueryRepo repository.ContainerProfileQueryRepo,
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	addContainerDto dto.AddContainer,
) error {
	err := CheckAccountQuota(
		accQueryRepo,
		addContainerDto.AccountId,
		containerProfileQueryRepo,
		*addContainerDto.ProfileId,
	)
	if err != nil {
		log.Printf("QuotaCheckError: %s", err)
		return err
	}

	_, err = containerQueryRepo.GetByHostname(addContainerDto.Hostname)
	if err == nil {
		log.Printf("ContainerHostnameAlreadyExists: %s", addContainerDto.Hostname)
		return errors.New("ContainerHostnameAlreadyExists")
	}

	containerId, err := containerCmdRepo.Add(addContainerDto)
	if err != nil {
		log.Printf("AddContainerError: %s", err)
		return errors.New("AddContainerInfraError")
	}

	err = accCmdRepo.UpdateQuotaUsage(addContainerDto.AccountId)
	if err != nil {
		log.Printf("UpdateAccountQuotaError: %s", err)
		return errors.New("UpdateAccountQuotaError")
	}

	log.Printf(
		"ContainerId '%s' (%s) created for AccountId '%s'.",
		containerId.String(),
		addContainerDto.ImageAddress.String(),
		addContainerDto.AccountId.String(),
	)

	if !addContainerDto.AutoCreateMappings {
		return nil
	}

	return AddMappingsWithContainerId(
		containerQueryRepo,
		mappingQueryRepo,
		mappingCmdRepo,
		containerId,
	)
}
