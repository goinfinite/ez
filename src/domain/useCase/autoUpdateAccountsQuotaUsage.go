package useCase

import (
	"log/slog"

	"github.com/speedianet/control/src/domain/repository"
)

func AutoUpdateAccountsQuotaUsage(
	accountQueryRepo repository.AccountQueryRepo,
	accountCmdRepo repository.AccountCmdRepo,
) {
	accounts, err := accountQueryRepo.Read()
	if err != nil {
		slog.Error("ReadAccountsInfraError", slog.Any("error", err))
		return
	}

	for _, account := range accounts {
		err := accountCmdRepo.UpdateQuotaUsage(account.Id)
		if err != nil {
			slog.Error(
				"UpdateQuotaUsageInfraError",
				slog.Uint64("accountId", account.Id.Uint64()),
				slog.Any("error", err),
			)
			continue
		}
	}
}
