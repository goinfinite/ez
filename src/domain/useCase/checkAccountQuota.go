package useCase

import (
	"errors"
	"log"

	"github.com/goinfinite/fleet/src/domain/repository"
	"github.com/goinfinite/fleet/src/domain/valueObject"
)

func CheckAccountQuota(
	accQueryRepo repository.AccQueryRepo,
	accId valueObject.AccountId,
	containerProfileQueryRepo repository.ContainerProfileQueryRepo,
	profileId valueObject.ContainerProfileId,
) error {
	accEntity, err := accQueryRepo.GetById(accId)
	if err != nil {
		log.Printf("GetAccountInfoError: %s", err)
		return errors.New("GetAccountInfoInfraError")
	}

	containerProfileEntity, err := containerProfileQueryRepo.GetById(profileId)
	if err != nil {
		log.Printf("GetContainerProfileError: %s", err)
		return errors.New("GetContainerProfileInfraError")
	}

	quotaCpu := accEntity.Quota.CpuCores
	quotaMemory := accEntity.Quota.MemoryBytes

	quotaUsageCpu := accEntity.QuotaUsage.CpuCores
	quotaUsageMemory := accEntity.QuotaUsage.MemoryBytes

	containerCpuLimit := containerProfileEntity.BaseSpecs.CpuCores
	containerMemoryLimit := containerProfileEntity.BaseSpecs.MemoryBytes

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
