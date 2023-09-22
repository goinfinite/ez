package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/sfm/src/domain/repository"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

func DeleteContainer(
	containerQueryRepo repository.ContainerQueryRepo,
	containerCmdRepo repository.ContainerCmdRepo,
	accCmdRepo repository.AccCmdRepo,
	accId valueObject.AccountId,
	containerId valueObject.ContainerId,
) error {
	_, err := containerQueryRepo.GetById(accId, containerId)
	if err != nil {
		return errors.New("ContainerNotFound")
	}

	err = containerCmdRepo.Delete(accId, containerId)
	if err != nil {
		return errors.New("DeleteContainerError")
	}

	log.Printf("ContainerId '%v' deleted.", containerId)

	err = accCmdRepo.UpdateQuotaUsage(accId)
	if err != nil {
		log.Printf("UpdateAccountQuotaError: %s", err)
		return errors.New("UpdateAccountQuotaError")
	}

	return nil
}
