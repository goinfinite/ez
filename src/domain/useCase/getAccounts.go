package useCase

import (
	"errors"
	"log"

	"github.com/goinfinite/fleet/src/domain/entity"
	"github.com/goinfinite/fleet/src/domain/repository"
)

func GetAccounts(
	accQueryRepo repository.AccQueryRepo,
) ([]entity.Account, error) {
	accs, err := accQueryRepo.Get()
	if err != nil {
		log.Printf("GetAccountsError: %s", err)
		return nil, errors.New("GetAccountsInfraError")
	}
	return accs, nil
}
