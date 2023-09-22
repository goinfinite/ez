package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/sfm/src/domain/repository"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

func AccQuotaCheck(
	accQueryRepo repository.AccQueryRepo,
	accId valueObject.AccountId,
	containerSpecs valueObject.ContainerSpecs,
) error {
	accEntity, err := accQueryRepo.GetById(accId)
	if err != nil {
		log.Printf("GetAccountInfoError: %s", err)
		return errors.New("GetAccountInfoInfraError")
	}

	accQuota := accEntity.Quota
	accQuotaUsage := accEntity.QuotaUsage

	if accQuotaUsage.CpuCores+containerSpecs.CpuCores > accQuota.CpuCores {
		log.Printf("CpuQuotaUsageExceeded: %s", err)
		return errors.New("CpuQuotaUsageExceeded")
	}

	if accQuotaUsage.MemoryBytes+containerSpecs.MemoryBytes > accQuota.MemoryBytes {
		log.Printf("MemoryQuotaUsageExceeded: %s", err)
		return errors.New("MemoryQuotaUsageExceeded")
	}

	return nil
}
