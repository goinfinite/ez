package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func CheckAccountQuota(
	accountQueryRepo repository.AccountQueryRepo,
	accId valueObject.AccountId,
	containerProfileQueryRepo repository.ContainerProfileQueryRepo,
	newProfileId valueObject.ContainerProfileId,
	prevProfileId *valueObject.ContainerProfileId,
) error {
	accEntity, err := accountQueryRepo.GetById(accId)
	if err != nil {
		log.Printf("GetAccountInfoError: %s", err)
		return errors.New("GetAccountInfoInfraError")
	}

	newProfile, err := containerProfileQueryRepo.GetById(newProfileId)
	if err != nil {
		log.Printf("GetNewContainerProfileError: %s", err)
		return errors.New("GetNewContainerProfileInfraError")
	}

	var prevProfilePtr *entity.ContainerProfile
	if prevProfileId != nil {
		prevProfile, err := containerProfileQueryRepo.GetById(*prevProfileId)
		if err != nil {
			log.Printf("GetPrevContainerProfileError: %s", err)
			return errors.New("GetPrevContainerProfileInfraError")
		}
		prevProfilePtr = &prevProfile
	}

	accCpuLimit := accEntity.Quota.CpuCores
	accMemoryLimit := accEntity.Quota.MemoryBytes

	accCpuUsage := accEntity.QuotaUsage.CpuCores
	accMemoryUsage := accEntity.QuotaUsage.MemoryBytes
	if prevProfilePtr != nil {
		accCpuUsage -= prevProfilePtr.BaseSpecs.CpuCores
		accMemoryUsage -= prevProfilePtr.BaseSpecs.MemoryBytes
	}

	newContainerCpuLimit := newProfile.BaseSpecs.CpuCores
	newContainerMemoryLimit := newProfile.BaseSpecs.MemoryBytes

	if accCpuUsage+newContainerCpuLimit > accCpuLimit {
		return errors.New("CpuQuotaUsageExceeded")
	}

	if accMemoryUsage+newContainerMemoryLimit > accMemoryLimit {
		return errors.New("MemoryQuotaUsageExceeded")
	}

	return nil
}
