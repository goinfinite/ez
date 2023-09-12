package useCase

import (
	"log"

	"github.com/speedianet/sfm/src/domain/repository"
)

func AutoUpdateAccountsQuotaUsage(
	accCmdRepo repository.AccCmdRepo,
) {
	err := accCmdRepo.UpdateQuotasUsage()
	if err != nil {
		log.Printf("UpdateQuotasUsageError: %v", err)
	}
}
