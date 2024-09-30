package dbModel

import (
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type ContainerProfile struct {
	ID                     uint64 `gorm:"primarykey"`
	AccountID              uint64 `gorm:"not null"`
	Name                   string `gorm:"not null"`
	BaseSpecs              string `gorm:"not null"`
	MaxSpecs               *string
	ScalingPolicy          *string
	ScalingThreshold       *uint
	ScalingMaxDurationSecs *uint
	ScalingIntervalSecs    *uint
	HostMinCapacityPercent *uint8
}

func (ContainerProfile) TableName() string {
	return "container_profiles"
}

func (model ContainerProfile) InitialEntries() []interface{} {
	initialProfilesEntities := entity.InitialContainerProfiles()

	initialProfiles := []interface{}{}
	for _, profileEntity := range initialProfilesEntities {
		profileModel, _ := model.ToModel(profileEntity)
		initialProfiles = append(initialProfiles, profileModel)
	}

	return initialProfiles
}

func NewContainerProfile(
	accountId uint64,
	name, baseSpecs string,
	maxSpecs, scalingPolicy *string,
	scalingThreshold, scalingMaxDurationSecs, scalingIntervalSecs *uint,
	hostMinCapacityPercent *uint8,
) ContainerProfile {
	return ContainerProfile{
		AccountID:              accountId,
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

func (ContainerProfile) ToModel(
	entity entity.ContainerProfile,
) (ContainerProfile, error) {
	var maxSpecsPtr *string
	if entity.MaxSpecs != nil {
		maxSpecs := entity.MaxSpecs.String()
		maxSpecsPtr = &maxSpecs
	}

	var scalingPolicyPtr *string
	if entity.ScalingPolicy != nil {
		scalingPolicy := entity.ScalingPolicy.String()
		scalingPolicyPtr = &scalingPolicy
	}

	var hostMinCapacityPercentPtr *uint8
	if entity.HostMinCapacityPercent != nil {
		hostMinCapacityPercentUint8 := entity.HostMinCapacityPercent.Uint8()
		hostMinCapacityPercentPtr = &hostMinCapacityPercentUint8
	}

	return ContainerProfile{
		ID:                     entity.Id.Uint64(),
		AccountID:              entity.AccountId.Uint64(),
		Name:                   entity.Name.String(),
		BaseSpecs:              entity.BaseSpecs.String(),
		MaxSpecs:               maxSpecsPtr,
		ScalingPolicy:          scalingPolicyPtr,
		ScalingThreshold:       entity.ScalingThreshold,
		ScalingMaxDurationSecs: entity.ScalingMaxDurationSecs,
		ScalingIntervalSecs:    entity.ScalingIntervalSecs,
		HostMinCapacityPercent: hostMinCapacityPercentPtr,
	}, nil
}

func (model ContainerProfile) ToEntity() (
	profileEntity entity.ContainerProfile, err error,
) {
	accountId, err := valueObject.NewAccountId(model.AccountID)
	if err != nil {
		return profileEntity, err
	}

	profileId, err := valueObject.NewContainerProfileId(model.ID)
	if err != nil {
		return profileEntity, err
	}

	name, err := valueObject.NewContainerProfileName(model.Name)
	if err != nil {
		return profileEntity, err
	}

	baseSpecs, err := valueObject.NewContainerSpecsFromString(model.BaseSpecs)
	if err != nil {
		return profileEntity, err
	}

	var maxSpecsPtr *valueObject.ContainerSpecs
	if model.MaxSpecs != nil {
		maxSpecs, err := valueObject.NewContainerSpecsFromString(*model.MaxSpecs)
		if err != nil {
			return profileEntity, err
		}
		maxSpecsPtr = &maxSpecs
	}

	var scalingPolicyPtr *valueObject.ScalingPolicy
	if model.ScalingPolicy != nil {
		scalingPolicy, err := valueObject.NewScalingPolicy(*model.ScalingPolicy)
		if err != nil {
			return profileEntity, err
		}
		scalingPolicyPtr = &scalingPolicy
	}

	var hostMinCapacityPercentPtr *valueObject.HostMinCapacity
	if model.HostMinCapacityPercent != nil {
		hostMinCapacityPercent, err := valueObject.NewHostMinCapacity(
			*model.HostMinCapacityPercent,
		)
		if err != nil {
			return profileEntity, err
		}
		hostMinCapacityPercentPtr = &hostMinCapacityPercent
	}

	return entity.NewContainerProfile(
		profileId, accountId, name, baseSpecs, maxSpecsPtr, scalingPolicyPtr,
		model.ScalingThreshold, model.ScalingMaxDurationSecs,
		model.ScalingIntervalSecs, hostMinCapacityPercentPtr,
	)
}
