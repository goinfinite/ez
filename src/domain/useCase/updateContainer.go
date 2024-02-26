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
	updateDto dto.UpdateContainer,
) error {
	containerEntity, err := containerQueryRepo.GetById(updateDto.ContainerId)
	if err != nil {
		return errors.New("ContainerNotFound")
	}

	shouldUpdateQuota := updateDto.ProfileId != nil
	if shouldUpdateQuota {
		err = CheckAccountQuota(
			accQueryRepo,
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
		err = accCmdRepo.UpdateQuotaUsage(updateDto.AccountId)
		if err != nil {
			log.Printf("UpdateAccountQuotaError: %s", err)
			return errors.New("UpdateAccountQuotaError")
		}
	}

	return nil
}
