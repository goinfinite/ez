package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/repository"
)

const AutoRefreshAccountQuotasTimeIntervalSecs uint16 = 900

func RefreshAccountQuotas(
	accountQueryRepo repository.AccountQueryRepo,
	accountCmdRepo repository.AccountCmdRepo,
) error {
	accountEntities, err := accountQueryRepo.Read()
	if err != nil {
		slog.Error("ReadAccountsInfraError", slog.Any("error", err))
		return errors.New("ReadAccountsInfraError")
	}

	for _, accountEntity := range accountEntities {
		err := accountCmdRepo.UpdateQuotaUsage(accountEntity.Id)
		if err != nil {
			slog.Debug(
				"UpdateQuotaUsageInfraError",
				slog.Uint64("accountId", accountEntity.Id.Uint64()),
				slog.Any("error", err),
			)
			continue
		}
	}

	return nil
}
