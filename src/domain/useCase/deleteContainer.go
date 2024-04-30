package useCase

import (
	"errors"
	"log"
	"time"

	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func mappingsJanitor(
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	containerId valueObject.ContainerId,
) error {
	targets, err := mappingQueryRepo.GetTargetsByContainerId(containerId)
	if err != nil {
		log.Printf("[%v] GetTargetsByContainerIdError: %s", containerId, err)
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

	mappings, err := mappingQueryRepo.Get()
	if err != nil {
		return nil
	}

	if len(mappings) == 0 {
		return nil
	}

	nowEpoch := time.Now().Unix()
	for _, mapping := range mappings {
		if len(mapping.Targets) != 0 {
			continue
		}

		isMappingTooRecent := nowEpoch-mapping.CreatedAt.Get() < 60
		if isMappingTooRecent {
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
	_, err := containerQueryRepo.GetById(containerId)
	if err != nil {
		log.Printf("ContainerNotFound: %s", err)
		return errors.New("ContainerNotFound")
	}

	err = mappingsJanitor(mappingQueryRepo, mappingCmdRepo, containerId)
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
