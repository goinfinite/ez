package useCase

import (
	"errors"
	"log"

	"github.com/goinfinite/fleet/src/domain/repository"
	"github.com/goinfinite/fleet/src/domain/valueObject"
)

func DeleteAccount(
	accQueryRepo repository.AccQueryRepo,
	accCmdRepo repository.AccCmdRepo,
	accountId valueObject.AccountId,
) error {
	_, err := accQueryRepo.GetById(accountId)
	if err != nil {
		return errors.New("AccountNotFound")
	}

	err = accCmdRepo.Delete(accountId)
	if err != nil {
		log.Printf("DeleteAccountError: %s", err)
		return errors.New("DeleteAccountInfraError")
	}

	log.Printf("AccountId '%v' deleted.", accountId)

	return nil
}
