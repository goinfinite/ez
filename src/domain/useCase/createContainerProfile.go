package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

var (
	ContainerProfileDefaultScalingThreshold       uint = 80
	ContainerProfileDefaultScalingMaxDurationSecs uint = 3600
	ContainerProfileDefaultScalingIntervalSecs    uint = 86400
)

func CreateContainerProfile(
	containerProfileCmdRepo repository.ContainerProfileCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateContainerProfile,
) error {
	if createDto.MaxSpecs != nil {
		if createDto.ScalingPolicy == nil {
			defaultPolicy := valueObject.DefaultScalingPolicy()
			createDto.ScalingPolicy = &defaultPolicy
		}

		if createDto.ScalingThreshold == nil {
			createDto.ScalingThreshold = &ContainerProfileDefaultScalingThreshold
		}

		if createDto.ScalingMaxDurationSecs == nil {
			createDto.ScalingMaxDurationSecs = &ContainerProfileDefaultScalingMaxDurationSecs
		}

		if createDto.ScalingIntervalSecs == nil {
			createDto.ScalingIntervalSecs = &ContainerProfileDefaultScalingIntervalSecs
		}

		if createDto.HostMinCapacityPercent == nil {
			defaultHostMinCapacity := valueObject.DefaultHostMinCapacity()
			createDto.HostMinCapacityPercent = &defaultHostMinCapacity
		}
	}

	profileId, err := containerProfileCmdRepo.Create(createDto)
	if err != nil {
		slog.Error("CreateContainerProfileInfraError", slog.Any("error", err))
		return errors.New("CreateContainerProfileInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateContainerProfile(createDto, profileId)

	return nil
}
