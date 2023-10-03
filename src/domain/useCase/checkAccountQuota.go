package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/sfm/src/domain/repository"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

func CheckAccountQuota(
	accQueryRepo repository.AccQueryRepo,
	accId valueObject.AccountId,
	resourceProfileQueryRepo repository.ResourceProfileQueryRepo,
	resourceProfileId valueObject.ResourceProfileId,
) error {
	accEntity, err := accQueryRepo.GetById(accId)
	if err != nil {
		log.Printf("GetAccountInfoError: %s", err)
		return errors.New("GetAccountInfoInfraError")
	}

	resourceProfileEntity, err := resourceProfileQueryRepo.GetById(resourceProfileId)
	if err != nil {
		log.Printf("GetResourceProfileError: %s", err)
		return errors.New("GetResourceProfileInfraError")
	}

	quotaCpu := accEntity.Quota.CpuCores
	quotaMemory := accEntity.Quota.MemoryBytes

	quotaUsageCpu := accEntity.QuotaUsage.CpuCores
	quotaUsageMemory := accEntity.QuotaUsage.MemoryBytes

	containerCpuLimit := resourceProfileEntity.BaseSpecs.CpuCores
	containerMemoryLimit := resourceProfileEntity.BaseSpecs.MemoryBytes

	if quotaUsageCpu+containerCpuLimit > quotaCpu {
		log.Printf("CpuQuotaUsageExceeded: %s", err)
		return errors.New("CpuQuotaUsageExceeded")
	}

	if quotaUsageMemory+containerMemoryLimit > quotaMemory {
		log.Printf("MemoryQuotaUsageExceeded: %s", err)
		return errors.New("MemoryQuotaUsageExceeded")
	}

	return nil
}
