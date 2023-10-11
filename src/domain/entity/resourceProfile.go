package entity

import (
	"github.com/speedianet/sfm/src/domain/valueObject"
)

type ResourceProfile struct {
	Id                     valueObject.ResourceProfileId   `json:"id"`
	Name                   valueObject.ResourceProfileName `json:"name"`
	BaseSpecs              valueObject.ContainerSpecs      `json:"baseSpecs"`
	MaxSpecs               *valueObject.ContainerSpecs     `json:"maxSpecs"`
	ScalingPolicy          *valueObject.ScalingPolicy      `json:"scalingPolicy"`
	ScalingThreshold       *uint64                         `json:"scalingThreshold"`
	ScalingMaxDurationSecs *uint64                         `json:"scalingMaxDurationSecs"`
	ScalingIntervalSecs    *uint64                         `json:"scalingIntervalSecs"`
	HostMinCapacityPercent *valueObject.HostMinCapacity    `json:"hostMinCapacityPercent"`
}

func NewResourceProfile(
	id valueObject.ResourceProfileId,
	name valueObject.ResourceProfileName,
	baseSpecs valueObject.ContainerSpecs,
	maxSpecs *valueObject.ContainerSpecs,
	scalingPolicy *valueObject.ScalingPolicy,
	scalingThreshold *uint64,
	scalingMaxDurationSecs *uint64,
	scalingIntervalSecs *uint64,
	hostMinCapacityPercent *valueObject.HostMinCapacity,
) (ResourceProfile, error) {
	return ResourceProfile{
		Id:                     id,
		Name:                   name,
		BaseSpecs:              baseSpecs,
		MaxSpecs:               maxSpecs,
		ScalingPolicy:          scalingPolicy,
		ScalingThreshold:       scalingThreshold,
		ScalingMaxDurationSecs: scalingMaxDurationSecs,
		ScalingIntervalSecs:    scalingIntervalSecs,
		HostMinCapacityPercent: hostMinCapacityPercent,
	}, nil
}
