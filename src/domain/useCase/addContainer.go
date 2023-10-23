package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/entity"
	"github.com/speedianet/sfm/src/domain/repository"
)

func AddContainer(
	containerCmdRepo repository.ContainerCmdRepo,
	accQueryRepo repository.AccQueryRepo,
	accCmdRepo repository.AccCmdRepo,
	containerProfileQueryRepo repository.ContainerProfileQueryRepo,
	addContainer dto.AddContainer,
) error {
	defaultContainerProfileId := entity.DefaultContainerProfile().Id
	if addContainer.ProfileId == nil {
		addContainer.ProfileId = &defaultContainerProfileId
	}

	err := CheckAccountQuota(
		accQueryRepo,
		addContainer.AccountId,
		containerProfileQueryRepo,
		*addContainer.ProfileId,
	)
	if err != nil {
		log.Printf("QuotaCheckError: %s", err)
		return err
	}

	err = containerCmdRepo.Add(addContainer)
	if err != nil {
		log.Printf("AddContainerError: %s", err)
		return errors.New("AddContainerInfraError")
	}

	err = accCmdRepo.UpdateQuotaUsage(addContainer.AccountId)
	if err != nil {
		log.Printf("UpdateAccountQuotaError: %s", err)
		return errors.New("UpdateAccountQuotaError")
	}

	return nil
}
