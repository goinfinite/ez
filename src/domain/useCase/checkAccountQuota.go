package useCase

import (
	"errors"
	"log/slog"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func CheckAccountQuota(
	accountQueryRepo repository.AccountQueryRepo,
	containerProfileQueryRepo repository.ContainerProfileQueryRepo,
	accountId valueObject.AccountId,
	newProfileId valueObject.ContainerProfileId,
	prevProfileId *valueObject.ContainerProfileId,
) error {
	accountEntity, err := accountQueryRepo.ReadById(accountId)
	if err != nil {
		slog.Error("GetAccountInfoInfraError", slog.Any("error", err))
		return errors.New("GetAccountInfoInfraError")
	}

	newProfileEntity, err := containerProfileQueryRepo.ReadById(newProfileId)
	if err != nil {
		slog.Error("GetNewContainerProfileInfraError", slog.Any("error", err))
		return errors.New("GetNewContainerProfileInfraError")
	}

	var prevProfileEntityPtr *entity.ContainerProfile
	if prevProfileId != nil {
		prevProfileEntity, err := containerProfileQueryRepo.ReadById(*prevProfileId)
		if err != nil {
			slog.Error("GetPrevContainerProfileInfraError", slog.Any("error", err))
			return errors.New("GetPrevContainerProfileInfraError")
		}
		prevProfileEntityPtr = &prevProfileEntity
	}

	accountCpuLimit := accountEntity.Quota.Millicores
	accountMemoryLimit := accountEntity.Quota.MemoryBytes

	accountCpuUsage := accountEntity.QuotaUsage.Millicores
	accountMemoryUsage := accountEntity.QuotaUsage.MemoryBytes
	if prevProfileEntityPtr != nil {
		accountCpuUsage -= prevProfileEntityPtr.BaseSpecs.Millicores
		accountMemoryUsage -= prevProfileEntityPtr.BaseSpecs.MemoryBytes
	}

	newContainerCpuLimit := newProfileEntity.BaseSpecs.Millicores
	newContainerMemoryLimit := newProfileEntity.BaseSpecs.MemoryBytes

	if accountCpuUsage+newContainerCpuLimit > accountCpuLimit {
		return errors.New("CpuQuotaUsageExceeded")
	}

	if accountMemoryUsage+newContainerMemoryLimit > accountMemoryLimit {
		return errors.New("MemoryQuotaUsageExceeded")
	}

	return nil
}
