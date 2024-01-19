package entity

import (
	"github.com/speedianet/control/src/domain/valueObject"
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

func InitialContainerProfiles() []ContainerProfile {
	defaultProfile := DefaultContainerProfile()

	profileId, _ := valueObject.NewContainerProfileId(2)
	profileName, _ := valueObject.NewContainerProfileName("small")
	baseSpecs, _ := valueObject.NewContainerSpecsFromString("1:2147483648")

	smallProfile := ContainerProfile{
		Id:        profileId,
		Name:      profileName,
		BaseSpecs: baseSpecs,
	}

	profileId, _ = valueObject.NewContainerProfileId(3)
	profileName, _ = valueObject.NewContainerProfileName("smallWithAutoScaling")
	baseSpecs, _ = valueObject.NewContainerSpecsFromString("1:2147483648")
	maxSpecs, _ := valueObject.NewContainerSpecsFromString("2:4294967296")
	scalingPolicy, _ := valueObject.NewScalingPolicy("cpu")
	scalingThreshold := uint64(80)
	scalingMaxDurationSecs := uint64(3600)
	scalingIntervalSecs := uint64(86400)
	hostMinCapacityPercent, _ := valueObject.NewHostMinCapacity(20)

	smallWithAutoScalingProfile, _ := NewContainerProfile(
		profileId,
		profileName,
		baseSpecs,
		&maxSpecs,
		&scalingPolicy,
		&scalingThreshold,
		&scalingMaxDurationSecs,
		&scalingIntervalSecs,
		&hostMinCapacityPercent,
	)

	profileId, _ = valueObject.NewContainerProfileId(4)
	profileName, _ = valueObject.NewContainerProfileName("medium")
	baseSpecs, _ = valueObject.NewContainerSpecsFromString("2:4294967296")

	mediumProfile := ContainerProfile{
		Id:        profileId,
		Name:      profileName,
		BaseSpecs: baseSpecs,
	}

	profileId, _ = valueObject.NewContainerProfileId(5)
	profileName, _ = valueObject.NewContainerProfileName("large")
	baseSpecs, _ = valueObject.NewContainerSpecsFromString("4:8589934592")

	largeProfile := ContainerProfile{
		Id:        profileId,
		Name:      profileName,
		BaseSpecs: baseSpecs,
	}

	return []ContainerProfile{
		defaultProfile,
		smallProfile,
		smallWithAutoScalingProfile,
		mediumProfile,
		largeProfile,
	}
}
