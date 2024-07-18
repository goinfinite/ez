package useCase

import (
	"errors"
	"log/slog"
	"time"

	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func mappingsJanitor(
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	containerId valueObject.ContainerId,
) error {
	targets, err := mappingQueryRepo.ReadTargetsByContainerId(containerId)
	if err != nil {
		slog.Error(
			"ReadTargetsByContainerIdInfraError",
			slog.String("containerId", containerId.String()),
			slog.Any("error", err),
		)
		return nil
	}

	for _, target := range targets {
		err = DeleteMappingTarget(
			mappingQueryRepo, mappingCmdRepo, target.MappingId, target.Id,
		)
		if err != nil {
			continue
		}
	}

	mappings, err := mappingQueryRepo.Read()
	if err != nil {
		return nil
	}

	if len(mappings) == 0 {
		return nil
	}

	nowEpoch := time.Now().Unix()
	for _, mapping := range mappings {
		if len(mapping.Targets) != 0 {
			continue
		}

		isMappingTooRecent := nowEpoch-mapping.CreatedAt.Read() < 60
		if isMappingTooRecent {
			continue
		}

		err = DeleteMapping(mappingQueryRepo, mappingCmdRepo, mapping.Id)
		if err != nil {
			continue
		}
	}

	return nil
}

func DeleteContainer(
	containerQueryRepo repository.ContainerQueryRepo,
	containerCmdRepo repository.ContainerCmdRepo,
	accountCmdRepo repository.AccountCmdRepo,
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	containerProxyCmdRepo repository.ContainerProxyCmdRepo,
	accountId valueObject.AccountId,
	containerId valueObject.ContainerId,
) error {
	_, err := containerQueryRepo.ReadById(containerId)
	if err != nil {
		return errors.New("ContainerNotFound")
	}

	err = mappingsJanitor(mappingQueryRepo, mappingCmdRepo, containerId)
	if err != nil {
		return err
	}

	err = containerProxyCmdRepo.Delete(containerId)
	if err != nil {
		slog.Error("DeleteContainerProxyInfraError", slog.Any("error", err))
		return errors.New("DeleteContainerProxyInfraError")
	}

	err = containerCmdRepo.Delete(accountId, containerId)
	if err != nil {
		slog.Error("DeleteContainerInfraError", slog.Any("error", err))
		return errors.New("DeleteContainerInfraError")
	}

	slog.Info("ContainerDeleted", slog.String("containerId", containerId.String()))

	err = accountCmdRepo.UpdateQuotaUsage(accountId)
	if err != nil {
		slog.Error("UpdateAccountQuotaInfraError", slog.Any("error", err))
		return errors.New("UpdateAccountQuotaInfraError")
	}

	return nil
}
