package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func mappingsJanitor(
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	containerEntity entity.Container,
) error {
	targets, err := mappingQueryRepo.GetTargetsByContainerId(containerEntity.Id)
	if err != nil {
		log.Printf("[%v] GetTargetsByContainerIdError: %s", containerEntity.Id, err)
		return nil
	}

	for _, target := range targets {
		err = mappingCmdRepo.DeleteTarget(target.Id)
		if err != nil {
			log.Printf("[%v] DeleteTargetError: %s", target.Id, err)
			continue
		}

		log.Printf("TargetId '%v' deleted.", target.Id)
	}

	mappings, err := mappingQueryRepo.GetByHostname(containerEntity.Hostname)
	if err != nil {
		return nil
	}

	if len(mappings) == 0 {
		return nil
	}

	for _, mapping := range mappings {
		if len(mapping.Targets) != 0 {
			continue
		}

		err = mappingCmdRepo.Delete(mapping.Id)
		if err != nil {
			log.Printf("[%v] DeleteMappingError: %s", mapping.Id, err)
			continue
		}

		log.Printf("MappingId '%v' deleted.", mapping.Id)
	}

	return nil
}

func DeleteContainer(
	containerQueryRepo repository.ContainerQueryRepo,
	containerCmdRepo repository.ContainerCmdRepo,
	accCmdRepo repository.AccCmdRepo,
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	accId valueObject.AccountId,
	containerId valueObject.ContainerId,
) error {
	containerEntity, err := containerQueryRepo.GetById(containerId)
	if err != nil {
		log.Printf("ContainerNotFound: %s", err)
		return errors.New("ContainerNotFound")
	}

	err = mappingsJanitor(mappingQueryRepo, mappingCmdRepo, containerEntity)
	if err != nil {
		return err
	}

	err = containerCmdRepo.Delete(accId, containerId)
	if err != nil {
		log.Printf("DeleteContainerError: %s", err)
		return errors.New("DeleteContainerInfraError")
	}

	log.Printf("ContainerId '%v' deleted.", containerId)

	err = accCmdRepo.UpdateQuotaUsage(accId)
	if err != nil {
		log.Printf("UpdateAccountQuotaError: %s", err)
		return errors.New("UpdateAccountQuotaError")
	}

	return nil
}
