package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/repository"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

func AddContainer(
	containerCmdRepo repository.ContainerCmdRepo,
	accQueryRepo repository.AccQueryRepo,
	accCmdRepo repository.AccCmdRepo,
	resourceProfileQueryRepo repository.ResourceProfileQueryRepo,
	addContainer dto.AddContainer,
) error {
	defaultResourceProfileId := valueObject.ResourceProfileId(1)

	if addContainer.ResourceProfileId == nil {
		addContainer.ResourceProfileId = &defaultResourceProfileId
	}

	err := CheckAccountQuota(
		accQueryRepo,
		addContainer.AccountId,
		resourceProfileQueryRepo,
		*addContainer.ResourceProfileId,
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
