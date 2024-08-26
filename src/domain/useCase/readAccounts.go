package useCase

import (
	"errors"
	"log/slog"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
)

func ReadAccounts(
	accountQueryRepo repository.AccountQueryRepo,
) ([]entity.Account, error) {
	accountsList, err := accountQueryRepo.Read()
	if err != nil {
		slog.Error("ReadAccountsInfraError", slog.Any("error", err))
		return accountsList, errors.New("ReadAccountsInfraError")
	}

	return accountsList, nil
}
