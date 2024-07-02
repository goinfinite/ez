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
	deleteDto dto.DeleteAccount,
) error {
	_, err := accountQueryRepo.ReadById(deleteDto.AccountId)
	if err != nil {
		return errors.New("AccountNotFound")
	}

	containers, err := containerQueryRepo.ReadByAccountId(deleteDto.AccountId)
	if err != nil {
		log.Printf("ReadContainersByAccIdError: %s", err)
		return errors.New("ReadContainersByAccIdInfraError")
	}

	if len(containers) > 0 {
		return errors.New("AccountHasContainers")
	}

	err = accountCmdRepo.Delete(deleteDto.AccountId)
	if err != nil {
		log.Printf("DeleteAccountError: %s", err)
		return errors.New("DeleteAccountInfraError")
	}

	eventType, _ := valueObject.NewSecurityEventType("account-deleted")
	createSecurityEventDto := dto.NewCreateSecurityEvent(
		eventType, nil, &deleteDto.IpAddress, &deleteDto.AccountId,
	)
	AsyncCreateSecurityEvent(securityCmdRepo, createSecurityEventDto)

	return nil
}
