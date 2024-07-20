package dto

import "github.com/speedianet/control/src/domain/valueObject"

type UpdateContainerProfile struct {
	Id                     valueObject.ContainerProfileId    `json:"id"`
	Name                   *valueObject.ContainerProfileName `json:"name"`
	BaseSpecs              *valueObject.ContainerSpecs       `json:"baseSpecs"`
	MaxSpecs               *valueObject.ContainerSpecs       `json:"maxSpecs"`
	ScalingPolicy          *valueObject.ScalingPolicy        `json:"scalingPolicy"`
	ScalingThreshold       *uint                             `json:"scalingThreshold"`
	ScalingMaxDurationSecs *uint                             `json:"scalingMaxDurationSecs"`
	ScalingIntervalSecs    *uint                             `json:"scalingIntervalSecs"`
	HostMinCapacityPercent *valueObject.HostMinCapacity      `json:"hostMinCapacityPercent"`
}

func NewUpdateContainerProfile(
	id valueObject.ContainerProfileId,
	name *valueObject.ContainerProfileName,
	baseSpecs *valueObject.ContainerSpecs,
	maxSpecs *valueObject.ContainerSpecs,
	scalingPolicy *valueObject.ScalingPolicy,
	scalingThreshold, scalingMaxDurationSecs, scalingIntervalSecs *uint,
	hostMinCapacityPercent *valueObject.HostMinCapacity,
) UpdateContainerProfile {
	return UpdateContainerProfile{
		Id:                     id,
		Name:                   name,
		BaseSpecs:              baseSpecs,
		MaxSpecs:               maxSpecs,
		ScalingPolicy:          scalingPolicy,
		ScalingThreshold:       scalingThreshold,
		ScalingMaxDurationSecs: scalingMaxDurationSecs,
		ScalingIntervalSecs:    scalingIntervalSecs,
		HostMinCapacityPercent: hostMinCapacityPercent,
	}
}
