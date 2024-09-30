package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
	"github.com/goinfinite/ez/src/domain/valueObject"
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
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
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

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateContainer(createDto, containerId)

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

	containerEntity, err := containerQueryRepo.ReadById(containerId)
	if err != nil {
		return errors.New("ContainerNotFound")
	}

	for _, portBinding := range containerEntity.PortBindings {
		createMappingDto := dto.NewCreateMapping(
			containerEntity.AccountId, &containerEntity.Hostname,
			portBinding.PublicPort, portBinding.Protocol,
			[]valueObject.ContainerId{containerId},
			createDto.OperatorAccountId, createDto.OperatorIpAddress,
		)
		err = CreateMapping(
			mappingQueryRepo, mappingCmdRepo, containerQueryRepo,
			activityRecordCmdRepo, createMappingDto,
		)
		if err != nil {
			continue
		}
	}

	return nil
}
