package useCase

import (
	"errors"
	"log/slog"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
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

	containerEntity, err := containerQueryRepo.ReadById(createDto.ContainerId)
	if err != nil {
		slog.Error("ReadContainerInfraError", slog.Any("error", err))
		return errors.New("ReadContainerInfraError")
	}

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
			slog.Uint64("publicPort", uint64(mappingEntity.PublicPort.Uint16())),
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
