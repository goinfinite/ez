package entity

import (
	"encoding/json"

	"github.com/goinfinite/ez/src/domain/valueObject"
)

type ContainerProfile struct {
	Id                        valueObject.ContainerProfileId   `json:"id"`
	AccountId                 valueObject.AccountId            `json:"accountId"`
	Name                      valueObject.ContainerProfileName `json:"name"`
	BaseSpecs                 valueObject.ContainerSpecs       `json:"baseSpecs"`
	MaxSpecs                  *valueObject.ContainerSpecs      `json:"maxSpecs"`
	ScalingPolicy             *valueObject.ScalingPolicy       `json:"scalingPolicy"`
	ScalingThreshold          *uint                            `json:"scalingThreshold"`
	ScalingMaxDurationSecs    *uint                            `json:"scalingMaxDurationSecs"`
	ScalingMaxDurationMinutes *uint16                          `json:"scalingMaxDurationMinutes"`
	ScalingMaxDurationHours   *uint8                           `json:"scalingMaxDurationHours"`
	ScalingIntervalSecs       *uint                            `json:"scalingIntervalSecs"`
	ScalingIntervalMinutes    *uint16                          `json:"scalingIntervalMinutes"`
	ScalingIntervalHours      *uint8                           `json:"scalingIntervalHours"`
	HostMinCapacityPercent    *valueObject.HostMinCapacity     `json:"hostMinCapacityPercent"`
}

func NewContainerProfile(
	id valueObject.ContainerProfileId,
	accountId valueObject.AccountId,
	name valueObject.ContainerProfileName,
	baseSpecs valueObject.ContainerSpecs,
	maxSpecs *valueObject.ContainerSpecs,
	scalingPolicy *valueObject.ScalingPolicy,
	scalingThreshold, scalingMaxDurationSecs, scalingIntervalSecs *uint,
	hostMinCapacityPercent *valueObject.HostMinCapacity,
) (ContainerProfile, error) {
	var scalingMaxDurationMinutesPtr *uint16
	var scalingMaxDurationHoursPtr *uint8
	if scalingMaxDurationSecs != nil {
		durationMinutes := uint16(*scalingMaxDurationSecs / 60)
		scalingMaxDurationMinutesPtr = &durationMinutes
		durationHours := uint8(*scalingMaxDurationSecs / 3600)
		scalingMaxDurationHoursPtr = &durationHours
	}

	var scalingIntervalMinutesPtr *uint16
	var scalingIntervalHoursPtr *uint8
	if scalingIntervalSecs != nil {
		intervalMinutes := uint16(*scalingIntervalSecs / 60)
		scalingIntervalMinutesPtr = &intervalMinutes
		intervalHours := uint8(*scalingIntervalSecs / 3600)
		scalingIntervalHoursPtr = &intervalHours
	}

	return ContainerProfile{
		Id:                        id,
		AccountId:                 accountId,
		Name:                      name,
		BaseSpecs:                 baseSpecs,
		MaxSpecs:                  maxSpecs,
		ScalingPolicy:             scalingPolicy,
		ScalingThreshold:          scalingThreshold,
		ScalingMaxDurationSecs:    scalingMaxDurationSecs,
		ScalingMaxDurationMinutes: scalingMaxDurationMinutesPtr,
		ScalingMaxDurationHours:   scalingMaxDurationHoursPtr,
		ScalingIntervalSecs:       scalingIntervalSecs,
		ScalingIntervalMinutes:    scalingIntervalMinutesPtr,
		ScalingIntervalHours:      scalingIntervalHoursPtr,
		HostMinCapacityPercent:    hostMinCapacityPercent,
	}, nil
}

func DefaultContainerProfile() ContainerProfile {
	profileId, _ := valueObject.NewContainerProfileId(1)
	profileName, _ := valueObject.NewContainerProfileName("default")
	baseSpecs, _ := valueObject.NewContainerSpecsFromString("500:1073741824:1")

	return ContainerProfile{
		Id:        profileId,
		AccountId: valueObject.SystemAccountId,
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
		AccountId: valueObject.SystemAccountId,
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
		profileId, valueObject.SystemAccountId, profileName, baseSpecs, &maxSpecs,
		&scalingPolicy, &scalingThreshold, &scalingMaxDurationSecs,
		&scalingIntervalSecs, &hostMinCapacityPercent,
	)

	profileId, _ = valueObject.NewContainerProfileId(4)
	profileName, _ = valueObject.NewContainerProfileName("medium")
	baseSpecs, _ = valueObject.NewContainerSpecsFromString("2000:4294967296:4")

	mediumProfile := ContainerProfile{
		Id:        profileId,
		AccountId: valueObject.SystemAccountId,
		Name:      profileName,
		BaseSpecs: baseSpecs,
	}

	profileId, _ = valueObject.NewContainerProfileId(5)
	profileName, _ = valueObject.NewContainerProfileName("large")
	baseSpecs, _ = valueObject.NewContainerSpecsFromString("4000:8589934592:5")

	largeProfile := ContainerProfile{
		Id:        profileId,
		AccountId: valueObject.SystemAccountId,
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
