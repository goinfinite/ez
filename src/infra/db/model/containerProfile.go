package dbModel

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type ContainerProfile struct {
	ID                     uint64 `gorm:"primarykey"`
	Name                   string `gorm:"not null"`
	BaseSpecs              string `gorm:"not null"`
	MaxSpecs               *string
	ScalingPolicy          *string
	ScalingThreshold       *uint
	ScalingMaxDurationSecs *uint
	ScalingIntervalSecs    *uint
	HostMinCapacityPercent *float64
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

	hostMinCapacityPercentFloat64 := entity.HostMinCapacityPercent.Float64()

	return ContainerProfile{
		ID:                     entity.Id.Uint64(),
		Name:                   entity.Name.String(),
		BaseSpecs:              entity.BaseSpecs.String(),
		MaxSpecs:               maxSpecsPtr,
		ScalingPolicy:          scalingPolicyPtr,
		ScalingThreshold:       entity.ScalingThreshold,
		ScalingMaxDurationSecs: entity.ScalingMaxDurationSecs,
		ScalingIntervalSecs:    entity.ScalingIntervalSecs,
		HostMinCapacityPercent: &hostMinCapacityPercentFloat64,
	}, nil
}

func (model ContainerProfile) ToEntity() (profileEntity entity.ContainerProfile, err error) {
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
		profileId, name, baseSpecs, maxSpecsPtr, scalingPolicyPtr, model.ScalingThreshold,
		model.ScalingMaxDurationSecs, model.ScalingIntervalSecs, hostMinCapacityPercentPtr,
	)
}

func (ContainerProfile) AddDtoToModel(
	dto dto.CreateContainerProfile,
) (ContainerProfile, error) {
	var maxSpecsPtr *string
	if dto.MaxSpecs != nil {
		maxSpecs := dto.MaxSpecs.String()
		maxSpecsPtr = &maxSpecs
	}

	var scalingPolicyPtr *string
	if dto.ScalingPolicy != nil {
		scalingPolicy := dto.ScalingPolicy.String()
		scalingPolicyPtr = &scalingPolicy
	}

	var hostMinCapacityPercentPtr *float64
	if dto.HostMinCapacityPercent != nil {
		hostMinCapacityPercentFloat64 := dto.HostMinCapacityPercent.Float64()
		hostMinCapacityPercentPtr = &hostMinCapacityPercentFloat64
	}

	return ContainerProfile{
		Name:                   dto.Name.String(),
		BaseSpecs:              dto.BaseSpecs.String(),
		MaxSpecs:               maxSpecsPtr,
		ScalingPolicy:          scalingPolicyPtr,
		ScalingThreshold:       dto.ScalingThreshold,
		ScalingMaxDurationSecs: dto.ScalingMaxDurationSecs,
		ScalingIntervalSecs:    dto.ScalingIntervalSecs,
		HostMinCapacityPercent: hostMinCapacityPercentPtr,
	}, nil
}
