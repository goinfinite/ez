package dto

import "github.com/speedianet/sfm/src/domain/valueObject"

type AddResourceProfile struct {
	Name             valueObject.ResourceProfileName `json:"name"`
	BaseSpecs        valueObject.ContainerSpecs      `json:"baseSpecs"`
	MaxSpecs         *valueObject.ContainerSpecs     `json:"maxSpecs"`
	ScalingPolicy    *valueObject.ScalingPolicy      `json:"scalingPolicy"`
	ScalingThreshold *uint64                         `json:"scalingThreshold"`
}

func NewAddResourceProfile(
	name valueObject.ResourceProfileName,
	baseSpecs valueObject.ContainerSpecs,
	maxSpecs *valueObject.ContainerSpecs,
	scalingPolicy *valueObject.ScalingPolicy,
	scalingThreshold *uint64,
) AddResourceProfile {
	return AddResourceProfile{
		Name:             name,
		BaseSpecs:        baseSpecs,
		MaxSpecs:         maxSpecs,
		ScalingPolicy:    scalingPolicy,
		ScalingThreshold: scalingThreshold,
	}
}
