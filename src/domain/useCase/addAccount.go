package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func AddAccount(
	accountQueryRepo repository.AccountQueryRepo,
	accountCmdRepo repository.AccountCmdRepo,
	addAccount dto.AddAccount,
) error {
	_, err := accountQueryRepo.GetByUsername(addAccount.Username)
	if err == nil {
		return errors.New("AccountAlreadyExists")
	}

	defaultQuota := valueObject.NewAccountQuotaWithDefaultValues()
	if addAccount.Quota == nil {
		addAccount.Quota = &defaultQuota
	}

	err = accountCmdRepo.Add(addAccount)
	if err != nil {
		log.Printf("AddAccountError: %s", err)
		return errors.New("AddAccountInfraError")
	}

	log.Printf("User '%v' added.", addAccount.Username.String())

	return nil
}
