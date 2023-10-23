package entity

import (
	"github.com/speedianet/sfm/src/domain/valueObject"
)

type ContainerProfile struct {
	Id                     valueObject.ContainerProfileId   `json:"id"`
	Name                   valueObject.ContainerProfileName `json:"name"`
	BaseSpecs              valueObject.ContainerSpecs       `json:"baseSpecs"`
	MaxSpecs               *valueObject.ContainerSpecs      `json:"maxSpecs"`
	ScalingPolicy          *valueObject.ScalingPolicy       `json:"scalingPolicy"`
	ScalingThreshold       *uint64                          `json:"scalingThreshold"`
	ScalingMaxDurationSecs *uint64                          `json:"scalingMaxDurationSecs"`
	ScalingIntervalSecs    *uint64                          `json:"scalingIntervalSecs"`
	HostMinCapacityPercent *valueObject.HostMinCapacity     `json:"hostMinCapacityPercent"`
}

func NewContainerProfile(
	id valueObject.ContainerProfileId,
	name valueObject.ContainerProfileName,
	baseSpecs valueObject.ContainerSpecs,
	maxSpecs *valueObject.ContainerSpecs,
	scalingPolicy *valueObject.ScalingPolicy,
	scalingThreshold *uint64,
	scalingMaxDurationSecs *uint64,
	scalingIntervalSecs *uint64,
	hostMinCapacityPercent *valueObject.HostMinCapacity,
) (ContainerProfile, error) {
	return ContainerProfile{
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

func DefaultContainerProfile() ContainerProfile {
	profileId, _ := valueObject.NewContainerProfileId(1)
	profileName, _ := valueObject.NewContainerProfileName("default")
	baseSpecs, _ := valueObject.NewContainerSpecsFromString("0.5:1073741824")

	return ContainerProfile{
		Id:                     profileId,
		Name:                   profileName,
		BaseSpecs:              baseSpecs,
		MaxSpecs:               nil,
		ScalingPolicy:          nil,
		ScalingThreshold:       nil,
		ScalingMaxDurationSecs: nil,
		ScalingIntervalSecs:    nil,
		HostMinCapacityPercent: nil,
	}
}
