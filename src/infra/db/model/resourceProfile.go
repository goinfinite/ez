package dbModel

import (
	"github.com/speedianet/sfm/src/domain/entity"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

type ResourceProfile struct {
	ID               uint   `gorm:"primarykey"`
	Name             string `gorm:"not null"`
	BaseSpecs        string `gorm:"not null"`
	MaxSpecs         *string
	ScalingPolicy    *string
	ScalingThreshold *uint64
}

func (ResourceProfile) TableName() string {
	return "resource_profiles"
}

func (ResourceProfile) ToModel(entity entity.ResourceProfile) (ResourceProfile, error) {
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

	var scalingThresholdPtr *uint64
	if entity.ScalingThreshold != nil {
		scalingThreshold := uint64(*entity.ScalingThreshold)
		scalingThresholdPtr = &scalingThreshold
	}

	return ResourceProfile{
		ID:               uint(entity.Id.Get()),
		Name:             entity.Name.String(),
		BaseSpecs:        entity.BaseSpecs.String(),
		MaxSpecs:         maxSpecsPtr,
		ScalingPolicy:    scalingPolicyPtr,
		ScalingThreshold: scalingThresholdPtr,
	}, nil
}

func (model ResourceProfile) ToEntity() (entity.ResourceProfile, error) {
	rpId, err := valueObject.NewResourceProfileId(model.ID)
	if err != nil {
		return entity.ResourceProfile{}, err
	}

	name, err := valueObject.NewResourceProfileName(model.Name)
	if err != nil {
		return entity.ResourceProfile{}, err
	}

	baseSpecs, err := valueObject.NewContainerSpecsFromString(model.BaseSpecs)
	if err != nil {
		return entity.ResourceProfile{}, err
	}

	var maxSpecsPtr *valueObject.ContainerSpecs
	if model.MaxSpecs != nil {
		maxSpecs, err := valueObject.NewContainerSpecsFromString(*model.MaxSpecs)
		if err != nil {
			return entity.ResourceProfile{}, err
		}
		maxSpecsPtr = &maxSpecs
	}

	var scalingPolicyPtr *valueObject.ScalingPolicy
	if model.ScalingPolicy != nil {
		scalingPolicy, err := valueObject.NewScalingPolicy(*model.ScalingPolicy)
		if err != nil {
			return entity.ResourceProfile{}, err
		}
		scalingPolicyPtr = &scalingPolicy
	}

	var scalingThresholdPtr *uint64
	if model.ScalingThreshold != nil {
		scalingThreshold := uint64(*model.ScalingThreshold)
		scalingThresholdPtr = &scalingThreshold
	}

	return entity.NewResourceProfile(
		rpId,
		name,
		baseSpecs,
		maxSpecsPtr,
		scalingPolicyPtr,
		scalingThresholdPtr,
	), nil
}
