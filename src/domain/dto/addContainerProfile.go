package dto

import "github.com/speedianet/control/src/domain/valueObject"

type AddContainerProfile struct {
	Name                   valueObject.ContainerProfileName `json:"name"`
	BaseSpecs              valueObject.ContainerSpecs       `json:"baseSpecs"`
	MaxSpecs               *valueObject.ContainerSpecs      `json:"maxSpecs"`
	ScalingPolicy          *valueObject.ScalingPolicy       `json:"scalingPolicy"`
	ScalingThreshold       *uint64                          `json:"scalingThreshold"`
	ScalingMaxDurationSecs *uint64                          `json:"scalingMaxDurationSecs"`
	ScalingIntervalSecs    *uint64                          `json:"scalingIntervalSecs"`
	HostMinCapacityPercent *valueObject.HostMinCapacity     `json:"hostMinCapacityPercent"`
}

func NewAddContainerProfile(
	name valueObject.ContainerProfileName,
	baseSpecs valueObject.ContainerSpecs,
	maxSpecs *valueObject.ContainerSpecs,
	scalingPolicy *valueObject.ScalingPolicy,
	scalingThreshold *uint64,
	scalingMaxDurationSecs *uint64,
	scalingIntervalSecs *uint64,
	hostMinCapacityPercent *valueObject.HostMinCapacity,
) AddContainerProfile {
	return AddContainerProfile{
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
