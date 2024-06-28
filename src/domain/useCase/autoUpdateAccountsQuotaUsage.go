package useCase

import (
	"log"

	"github.com/speedianet/control/src/domain/repository"
)

func AutoUpdateAccountsQuotaUsage(
	accountQueryRepo repository.AccountQueryRepo,
	accountCmdRepo repository.AccountCmdRepo,
) {
	accs, err := accountQueryRepo.Get()
	if err != nil {
		log.Printf("GetAccountsError: %v", err)
		return
	}

	for _, acc := range accs {
		err := accountCmdRepo.UpdateQuotaUsage(acc.Id)
		if err != nil {
			log.Printf("UpdateQuotaUsageError: %v [accId: %s]", acc.Id, err)
			continue
		}
	}
}
