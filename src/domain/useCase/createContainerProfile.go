package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
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
			defaultThreshold := uint64(80)
			dto.ScalingThreshold = &defaultThreshold
		}

		if dto.ScalingMaxDurationSecs == nil {
			defaultMaxDuration := uint64(3600)
			dto.ScalingMaxDurationSecs = &defaultMaxDuration
		}

		if dto.ScalingIntervalSecs == nil {
			defaultInterval := uint64(180)
			dto.ScalingIntervalSecs = &defaultInterval
		}

		if dto.HostMinCapacityPercent == nil {
			defaultHostMinCapacity := valueObject.DefaultHostMinCapacity()
			dto.HostMinCapacityPercent = &defaultHostMinCapacity
		}
	}

	err := containerProfileCmdRepo.Create(dto)
	if err != nil {
		log.Printf("CreateContainerProfileError: %v", err)
		return errors.New("CreateContainerProfileInfraError")
	}

	log.Printf("ContainerProfile '%s' added.", dto.Name.String())

	return nil
}
