package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
)

func UpdateContainer(
	containerQueryRepo repository.ContainerQueryRepo,
	containerCmdRepo repository.ContainerCmdRepo,
	accQueryRepo repository.AccQueryRepo,
	accCmdRepo repository.AccCmdRepo,
	containerProfileQueryRepo repository.ContainerProfileQueryRepo,
	updateContainer dto.UpdateContainer,
) error {
	currentContainer, err := containerQueryRepo.GetById(
		updateContainer.AccountId,
		updateContainer.ContainerId,
	)
	if err != nil {
		return errors.New("ContainerNotFound")
	}

	shouldUpdateQuota := updateContainer.ProfileId != nil
	if shouldUpdateQuota {
		err = CheckAccountQuota(
			accQueryRepo,
			updateContainer.AccountId,
			containerProfileQueryRepo,
			*updateContainer.ProfileId,
		)
		if err != nil {
			return err
		}
	}

	// Current OCI implementations does not support permanent container resources update.
	// Therefore, when updating container status, we also need to update the container
	// profile to guarantee that the container resources are up-to-date.
	if updateContainer.ProfileId == nil {
		updateContainer.ProfileId = &currentContainer.ProfileId
	}

	err = containerCmdRepo.Update(currentContainer, updateContainer)
	if err != nil {
		log.Printf("UpdateContainerError: %s", err)
		return errors.New("UpdateContainerInfraError")
	}

	if shouldUpdateQuota {
		err = accCmdRepo.UpdateQuotaUsage(updateContainer.AccountId)
		if err != nil {
			log.Printf("UpdateAccountQuotaError: %s", err)
			return errors.New("UpdateAccountQuotaError")
		}
	}

	return nil
}
