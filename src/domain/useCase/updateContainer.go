package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/repository"
)

func UpdateContainer(
	containerQueryRepo repository.ContainerQueryRepo,
	containerCmdRepo repository.ContainerCmdRepo,
	accQueryRepo repository.AccQueryRepo,
	accCmdRepo repository.AccCmdRepo,
	resourceProfileQueryRepo repository.ResourceProfileQueryRepo,
	updateContainer dto.UpdateContainer,
) error {
	_, err := containerQueryRepo.GetById(
		updateContainer.AccountId,
		updateContainer.ContainerId,
	)
	if err != nil {
		return errors.New("ContainerNotFound")
	}

	shouldUpdateQuota := updateContainer.ResourceProfileId != nil
	if shouldUpdateQuota {
		err = CheckAccountQuota(
			accQueryRepo,
			updateContainer.AccountId,
			resourceProfileQueryRepo,
			*updateContainer.ResourceProfileId,
		)
		if err != nil {
			return err
		}
	}

	err = containerCmdRepo.Update(updateContainer)
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
