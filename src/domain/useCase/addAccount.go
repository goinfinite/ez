package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/repository"
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
		return errors.New("AddAccountError")
	}

	log.Printf("User '%v' added.", addAccount.Username.String())

	return nil
}
