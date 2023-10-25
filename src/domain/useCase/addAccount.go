package useCase

import (
	"errors"
	"log"

	"github.com/goinfinite/fleet/src/domain/dto"
	"github.com/goinfinite/fleet/src/domain/repository"
)

func AddAccount(
	accQueryRepo repository.AccQueryRepo,
	accCmdRepo repository.AccCmdRepo,
	addAccount dto.AddAccount,
) error {
	_, err := accQueryRepo.GetByUsername(addAccount.Username)
	if err == nil {
		return errors.New("AccountAlreadyExists")
	}

	err = accCmdRepo.Add(addAccount)
	if err != nil {
		log.Printf("AddAccountError: %s", err)
		return errors.New("AddAccountInfraError")
	}

	log.Printf("User '%v' added.", addAccount.Username.String())

	return nil
}
