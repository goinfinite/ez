package entity

import (
	"encoding/json"

	"github.com/speedianet/control/src/domain/valueObject"
)

type ContainerProfile struct {
	Id                     valueObject.ContainerProfileId   `json:"id"`
	Name                   valueObject.ContainerProfileName `json:"name"`
	BaseSpecs              valueObject.ContainerSpecs       `json:"baseSpecs"`
	MaxSpecs               *valueObject.ContainerSpecs      `json:"maxSpecs"`
	ScalingPolicy          *valueObject.ScalingPolicy       `json:"scalingPolicy"`
	ScalingThreshold       *uint                            `json:"scalingThreshold"`
	ScalingMaxDurationSecs *uint                            `json:"scalingMaxDurationSecs"`
	ScalingIntervalSecs    *uint                            `json:"scalingIntervalSecs"`
	HostMinCapacityPercent *valueObject.HostMinCapacity     `json:"hostMinCapacityPercent"`
}

func NewContainerProfile(
	id valueObject.ContainerProfileId,
	name valueObject.ContainerProfileName,
	baseSpecs valueObject.ContainerSpecs,
	maxSpecs *valueObject.ContainerSpecs,
	scalingPolicy *valueObject.ScalingPolicy,
	scalingThreshold, scalingMaxDurationSecs, scalingIntervalSecs *uint,
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
	baseSpecs, _ := valueObject.NewContainerSpecsFromString("500:1073741824:1")

	return ContainerProfile{
		Id:        profileId,
		Name:      profileName,
		BaseSpecs: baseSpecs,
	}
}

func InitialContainerProfiles() []ContainerProfile {
	defaultProfile := DefaultContainerProfile()

	profileId, _ := valueObject.NewContainerProfileId(2)
	profileName, _ := valueObject.NewContainerProfileName("small")
	baseSpecs, _ := valueObject.NewContainerSpecsFromString("1000:2147483648:2")

	smallProfile := ContainerProfile{
		Id:        profileId,
		Name:      profileName,
		BaseSpecs: baseSpecs,
	}

	profileId, _ = valueObject.NewContainerProfileId(3)
	profileName, _ = valueObject.NewContainerProfileName("smallWithAutoScaling")
	baseSpecs, _ = valueObject.NewContainerSpecsFromString("1000:2147483648:2")
	maxSpecs, _ := valueObject.NewContainerSpecsFromString("2000:4294967296:3")
	scalingPolicy, _ := valueObject.NewScalingPolicy("cpu")
	scalingThreshold := uint(80)
	scalingMaxDurationSecs := uint(3600)
	scalingIntervalSecs := uint(86400)
	hostMinCapacityPercent, _ := valueObject.NewHostMinCapacity(20)

	smallWithAutoScalingProfile, _ := NewContainerProfile(
		profileId, profileName, baseSpecs, &maxSpecs, &scalingPolicy, &scalingThreshold,
		&scalingMaxDurationSecs, &scalingIntervalSecs, &hostMinCapacityPercent,
	)

	profileId, _ = valueObject.NewContainerProfileId(4)
	profileName, _ = valueObject.NewContainerProfileName("medium")
	baseSpecs, _ = valueObject.NewContainerSpecsFromString("2000:4294967296:4")

	mediumProfile := ContainerProfile{
		Id:        profileId,
		Name:      profileName,
		BaseSpecs: baseSpecs,
	}

	profileId, _ = valueObject.NewContainerProfileId(5)
	profileName, _ = valueObject.NewContainerProfileName("large")
	baseSpecs, _ = valueObject.NewContainerSpecsFromString("4000:8589934592:5")

	largeProfile := ContainerProfile{
		Id:        profileId,
		Name:      profileName,
		BaseSpecs: baseSpecs,
	}

	return []ContainerProfile{
		defaultProfile, smallProfile, smallWithAutoScalingProfile,
		mediumProfile, largeProfile,
	}
}

func (entity ContainerProfile) JsonSerialize() string {
	if entity.MaxSpecs == nil {
		entity.MaxSpecs = &valueObject.ContainerSpecs{}
	}

	jsonBytes, _ := json.Marshal(entity)
	return string(jsonBytes)
}
