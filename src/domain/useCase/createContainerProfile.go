package useCase

import (
	"errors"
	"log/slog"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

var (
	ContainerProfileDefaultScalingThreshold       uint = 80
	ContainerProfileDefaultScalingMaxDurationSecs uint = 3600
	ContainerProfileDefaultScalingIntervalSecs    uint = 180
)

func CreateContainerProfile(
	containerProfileCmdRepo repository.ContainerProfileCmdRepo,
	dto dto.CreateContainerProfile,
) error {
	if dto.MaxSpecs != nil {
		if dto.ScalingPolicy == nil {
			defaultPolicy := valueObject.DefaultScalingPolicy()
			dto.ScalingPolicy = &defaultPolicy
		}

		if dto.ScalingThreshold == nil {
			dto.ScalingThreshold = &ContainerProfileDefaultScalingThreshold
		}

		if dto.ScalingMaxDurationSecs == nil {
			dto.ScalingMaxDurationSecs = &ContainerProfileDefaultScalingMaxDurationSecs
		}

		if dto.ScalingIntervalSecs == nil {
			dto.ScalingIntervalSecs = &ContainerProfileDefaultScalingIntervalSecs
		}

		if dto.HostMinCapacityPercent == nil {
			defaultHostMinCapacity := valueObject.DefaultHostMinCapacity()
			dto.HostMinCapacityPercent = &defaultHostMinCapacity
		}
	}

	err := containerProfileCmdRepo.Create(dto)
	if err != nil {
		slog.Error(
			"CreateContainerProfileInfraError",
			slog.Any("error", err),
		)
		return errors.New("CreateContainerProfileInfraError")
	}

	slog.Info(
		"ContainerProfileCreated",
		slog.String("name", dto.Name.String()),
	)

	return nil
}
