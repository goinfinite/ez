package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
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
		"ContainerId '%s' with image address '%s' created for AccountId '%s'.",
		containerId.String(),
		addContainerDto.ImageAddr.String(),
		addContainerDto.AccountId.String(),
	)

	if !addContainerDto.AutoCreateMappings {
		return nil
	}

	for _, portBinding := range addContainerDto.PortBindings {
		addMappingDto := dto.NewAddMapping(
			addContainerDto.AccountId,
			&addContainerDto.Hostname,
			portBinding.PublicPort,
			portBinding.Protocol,
			[]valueObject.ContainerId{containerId},
		)
		err = AddMapping(
			mappingQueryRepo,
			mappingCmdRepo,
			containerQueryRepo,
			addMappingDto,
		)
		if err != nil {
			log.Printf("AddMappingError: %s", err)
			return errors.New("AddMappingError")
		}
	}

	return nil
}
