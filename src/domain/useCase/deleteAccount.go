package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func DeleteAccount(
	accountQueryRepo repository.AccountQueryRepo,
	accountCmdRepo repository.AccountCmdRepo,
	accountId valueObject.AccountId,
	containerQueryRepo repository.ContainerQueryRepo,
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

	log.Printf("AccountId '%v' deleted.", accountId)

	return nil
}
