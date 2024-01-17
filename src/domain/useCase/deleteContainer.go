package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func DeleteContainer(
	containerQueryRepo repository.ContainerQueryRepo,
	containerCmdRepo repository.ContainerCmdRepo,
	accCmdRepo repository.AccCmdRepo,
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	accId valueObject.AccountId,
	containerId valueObject.ContainerId,
) error {
	_, err := containerQueryRepo.GetById(containerId)
	if err != nil {
		log.Printf("ContainerNotFound: %s", err)
		return errors.New("ContainerNotFound")
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

	targets, err := mappingQueryRepo.GetTargetsByContainerId(containerId)
	if err != nil {
		log.Printf("GetTargetByContainerIdError: %s", err)
		return errors.New("GetTargetByContainerIdInfraError")
	}

	for _, target := range targets {
		err = mappingCmdRepo.DeleteTarget(target.Id)
		if err != nil {
			log.Printf("DeleteTargetError: %s", err)
			return errors.New("DeleteTargetInfraError")
		}

		log.Printf("Mapping target '%v' deleted.", target.Id)
	}

	return nil
}
