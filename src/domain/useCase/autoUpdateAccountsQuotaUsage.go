package useCase

import (
	"log"

	"github.com/goinfinite/fleet/src/domain/repository"
)

func AutoUpdateAccountsQuotaUsage(
	accQueryRepo repository.AccQueryRepo,
	accCmdRepo repository.AccCmdRepo,
) {
	accs, err := accQueryRepo.Get()
	if err != nil {
		log.Printf("GetAccountsError: %v", err)
		return
	}

	for _, acc := range accs {
		err := accCmdRepo.UpdateQuotaUsage(acc.Id)
		if err != nil {
			log.Printf("UpdateQuotaUsageError: %v [accId: %s]", acc.Id, err)
			continue
		}
	}
}
