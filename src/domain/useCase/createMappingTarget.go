package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
)

func CreateMappingTarget(
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	containerQueryRepo repository.ContainerQueryRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateMappingTarget,
) error {
	mappingEntity, err := mappingQueryRepo.ReadById(createDto.MappingId)
	if err != nil {
		slog.Error("ReadMappingInfraError", slog.Any("error", err))
		return errors.New("ReadMappingInfraError")
	}

	readContainersDto := dto.ReadContainersRequest{
		Pagination:  ContainersDefaultPagination,
		ContainerId: &createDto.ContainerId,
	}

	responseDto, err := ReadContainers(containerQueryRepo, readContainersDto)
	if err != nil || len(responseDto.Containers) == 0 {
		return errors.New("ContainerNotFound")
	}
	containerEntity := responseDto.Containers[0]

	publicPortMatches := false
	for _, portBinding := range containerEntity.PortBindings {
		if portBinding.PublicPort != mappingEntity.PublicPort {
			continue
		}
		publicPortMatches = true
	}

	if !publicPortMatches {
		slog.Error(
			"ContainerDoesNotBindToMappingPublicPort",
			slog.String("containerId", createDto.ContainerId.String()),
			slog.String("publicPort", mappingEntity.PublicPort.String()),
		)
		return errors.New("ContainerDoesNotBindToMappingPublicPort")
	}

	for _, target := range mappingEntity.Targets {
		if target.ContainerId != createDto.ContainerId {
			continue
		}

		slog.Debug(
			"SkipExistingMappingTarget",
			slog.String("containerId", createDto.ContainerId.String()),
			slog.String("mappingId", createDto.MappingId.String()),
		)
		return nil
	}

	targetId, err := mappingCmdRepo.CreateTarget(createDto)
	if err != nil {
		slog.Error("CreateMappingTargetInfraError", slog.Any("error", err))
		return errors.New("CreateMappingTargetInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateMappingTarget(createDto, targetId)

	return nil
}
