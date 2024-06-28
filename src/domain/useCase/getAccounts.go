package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
)

func GetAccounts(
	accountQueryRepo repository.AccountQueryRepo,
) ([]entity.Account, error) {
	accs, err := accountQueryRepo.Get()
	if err != nil {
		log.Printf("GetAccountsError: %s", err)
		return nil, errors.New("GetAccountsInfraError")
	}
	return accs, nil
}
