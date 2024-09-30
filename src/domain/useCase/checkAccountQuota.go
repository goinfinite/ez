package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/repository"
	"github.com/goinfinite/ez/src/domain/valueObject"
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
		slog.Error("ReadAccountInfoInfraError", slog.Any("error", err))
		return errors.New("ReadAccountInfoInfraError")
	}

	newProfileEntity, err := containerProfileQueryRepo.ReadById(newProfileId)
	if err != nil {
		slog.Error("ReadNewContainerProfileInfraError", slog.Any("error", err))
		return errors.New("ReadNewContainerProfileInfraError")
	}

	var prevProfileEntityPtr *entity.ContainerProfile
	if prevProfileId != nil {
		prevProfileEntity, err := containerProfileQueryRepo.ReadById(*prevProfileId)
		if err != nil {
			slog.Error("ReadPrevContainerProfileInfraError", slog.Any("error", err))
			return errors.New("ReadPrevContainerProfileInfraError")
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
