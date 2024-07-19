package useCase

import (
	"errors"
	"log/slog"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
)

func CreateContainer(
	containerQueryRepo repository.ContainerQueryRepo,
	containerCmdRepo repository.ContainerCmdRepo,
	accountQueryRepo repository.AccountQueryRepo,
	accountCmdRepo repository.AccountCmdRepo,
	containerProfileQueryRepo repository.ContainerProfileQueryRepo,
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	containerProxyCmdRepo repository.ContainerProxyCmdRepo,
	createDto dto.CreateContainer,
) error {
	err := CheckAccountQuota(
		accountQueryRepo, containerProfileQueryRepo, createDto.AccountId,
		*createDto.ProfileId, nil,
	)
	if err != nil {
		return err
	}

	_, err = containerQueryRepo.ReadByHostname(createDto.Hostname)
	if err == nil {
		return errors.New("ContainerHostnameAlreadyExists")
	}

	containerId, err := containerCmdRepo.Create(createDto)
	if err != nil {
		slog.Error("CreateContainerInfraError", slog.Any("error", err))
		return errors.New("CreateContainerInfraError")
	}

	err = accountCmdRepo.UpdateQuotaUsage(createDto.AccountId)
	if err != nil {
		slog.Error("UpdateAccountQuotaInfraError", slog.Any("error", err))
		return errors.New("UpdateAccountQuotaInfraError")
	}

	slog.Info(
		"ContainerCreated",
		slog.String("containerId", containerId.String()),
		slog.String("imageAddress", createDto.ImageAddress.String()),
		slog.String("accountId", createDto.AccountId.String()),
	)

	if createDto.ImageAddress.IsSpeediaOs() {
		err = containerProxyCmdRepo.Create(containerId)
		if err != nil {
			slog.Error("CreateContainerProxyInfraError", slog.Any("error", err))
			return errors.New("CreateContainerProxyInfraError")
		}
	}

	if !createDto.AutoCreateMappings {
		return nil
	}

	err = CreateMappingsWithContainerId(
		containerQueryRepo, mappingQueryRepo, mappingCmdRepo,
		containerProxyCmdRepo, containerId,
	)
	if err != nil {
		return errors.New("CreateAutoMappingsError: " + err.Error())
	}

	return nil
}
