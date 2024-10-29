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
	containerImageQueryRepo repository.ContainerImageQueryRepo,
	containerImageCmdRepo repository.ContainerImageCmdRepo,
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

	readContainersDto := dto.ReadContainersRequest{
		Pagination:        ContainersDefaultPagination,
		ContainerHostname: &createDto.Hostname,
	}

	_, err = ReadContainers(containerQueryRepo, readContainersDto)
	if err == nil {
		return errors.New("ContainerHostnameAlreadyExists")
	}

	isInfiniteOs := createDto.ImageAddress.IsInfiniteOs()
	if createDto.ExistingContainerId != nil {
		readContainersDto = dto.ReadContainersRequest{
			Pagination:  ContainersDefaultPagination,
			ContainerId: createDto.ExistingContainerId,
		}

		responseDto, err := ReadContainers(containerQueryRepo, readContainersDto)
		if err != nil || len(responseDto.Containers) == 0 {
			return errors.New("ExistingContainerNotFound")
		}
		existingContainerEntity := responseDto.Containers[0]

		isInfiniteOs = existingContainerEntity.ImageAddress.IsInfiniteOs()

		createSnapshotDto := dto.NewCreateContainerSnapshotImage(
			*createDto.ExistingContainerId, nil, nil, nil, createDto.OperatorAccountId,
			createDto.OperatorIpAddress,
		)

		imageId, err := containerImageCmdRepo.CreateSnapshot(createSnapshotDto)
		if err != nil {
			slog.Error("CreateContainerSnapshotImageInfraError", slog.Any("error", err))
			return errors.New("CreateContainerSnapshotImageInfraError")
		}

		createDto.ImageId = &imageId
	}

	if createDto.ImageId != nil {
		imageEntity, err := containerImageQueryRepo.ReadById(
			createDto.AccountId, *createDto.ImageId,
		)
		if err != nil {
			return errors.New("ContainerImageNotFound")
		}

		createDto.ImageAddress = imageEntity.ImageAddress
		isInfiniteOs = imageEntity.ImageAddress.IsInfiniteOs()

		createDto.PortBindings = imageEntity.PortBindings
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

	if isInfiniteOs {
		err = containerProxyCmdRepo.Create(containerId)
		if err != nil {
			slog.Error("CreateContainerProxyInfraError", slog.Any("error", err))
			return errors.New("CreateContainerProxyInfraError")
		}
	}

	if !createDto.AutoCreateMappings {
		return nil
	}

	readContainersDto = dto.ReadContainersRequest{
		Pagination:  ContainersDefaultPagination,
		ContainerId: &containerId,
	}

	responseDto, err := ReadContainers(containerQueryRepo, readContainersDto)
	if err != nil || len(responseDto.Containers) == 0 {
		return errors.New("ContainerNotFound")
	}
	containerEntity := responseDto.Containers[0]

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
