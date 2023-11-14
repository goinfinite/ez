package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
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

	defaultQuota := valueObject.NewAccountQuotaWithDefaultValues()
	if addAccount.Quota == nil {
		addAccount.Quota = &defaultQuota
	}

	err = accCmdRepo.Add(addAccount)
	if err != nil {
		log.Printf("AddAccountError: %s", err)
		return errors.New("AddAccountInfraError")
	}

	log.Printf("User '%v' added.", addAccount.Username.String())

	return nil
}
