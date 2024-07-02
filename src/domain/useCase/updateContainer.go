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
	accountQueryRepo repository.AccountQueryRepo,
	accountCmdRepo repository.AccountCmdRepo,
	containerProfileQueryRepo repository.ContainerProfileQueryRepo,
	updateDto dto.UpdateContainer,
) error {
	containerEntity, err := containerQueryRepo.ReadById(updateDto.ContainerId)
	if err != nil {
		return errors.New("ContainerNotFound")
	}

	shouldUpdateQuota := updateDto.ProfileId != nil
	if shouldUpdateQuota {
		err = CheckAccountQuota(
			accountQueryRepo,
			updateDto.AccountId,
			containerProfileQueryRepo,
			*updateDto.ProfileId,
			&containerEntity.ProfileId,
		)
		if err != nil {
			return err
		}
	}

	err = containerCmdRepo.Update(updateDto)
	if err != nil {
		log.Printf("UpdateContainerError: %s", err)
		return errors.New("UpdateContainerInfraError")
	}

	if shouldUpdateQuota {
		err = accountCmdRepo.UpdateQuotaUsage(updateDto.AccountId)
		if err != nil {
			log.Printf("UpdateAccountQuotaError: %s", err)
			return errors.New("UpdateAccountQuotaError")
		}
	}

	return nil
}
