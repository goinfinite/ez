package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func DeleteAccount(
	accountQueryRepo repository.AccountQueryRepo,
	accountCmdRepo repository.AccountCmdRepo,
	containerQueryRepo repository.ContainerQueryRepo,
	securityCmdRepo repository.SecurityCmdRepo,
	accountId valueObject.AccountId,
	ipAddress *valueObject.IpAddress,
) error {
	_, err := accountQueryRepo.GetById(accountId)
	if err != nil {
		return errors.New("AccountNotFound")
	}

	containers, err := containerQueryRepo.GetByAccId(accountId)
	if err != nil {
		log.Printf("GetContainersByAccIdError: %s", err)
		return errors.New("GetContainersByAccIdInfraError")
	}

	if len(containers) > 0 {
		return errors.New("AccountHasContainers")
	}

	err = accountCmdRepo.Delete(accountId)
	if err != nil {
		log.Printf("DeleteAccountError: %s", err)
		return errors.New("DeleteAccountInfraError")
	}

	eventType, _ := valueObject.NewSecurityEventType("account-deleted")
	createSecurityEventDto := dto.NewCreateSecurityEvent(
		eventType, nil, ipAddress, &accountId,
	)
	err = CreateSecurityEvent(securityCmdRepo, createSecurityEventDto)
	if err != nil {
		return err
	}

	return nil
}
