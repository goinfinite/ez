package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
)

func ReadAccounts(
	accountQueryRepo repository.AccountQueryRepo,
) ([]entity.Account, error) {
	accountsList, err := accountQueryRepo.Read()
	if err != nil {
		log.Printf("ReadAccountsError: %s", err)
		return nil, errors.New("ReadAccountsInfraError")
	}

	return accountsList, nil
}
