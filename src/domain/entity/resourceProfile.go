package entity

import "github.com/speedianet/sfm/src/domain/valueObject"

type ResourceProfile struct {
	Id               valueObject.ResourceProfileId   `json:"id"`
	Name             valueObject.ResourceProfileName `json:"name"`
	BaseSpecs        valueObject.ContainerSpecs      `json:"baseSpecs"`
	MaxSpecs         *valueObject.ContainerSpecs     `json:"maxSpecs"`
	ScalingPolicy    *valueObject.ScalingPolicy      `json:"scalingPolicy"`
	ScalingThreshold *uint64                         `json:"scalingThreshold"`
}

func NewResourceProfile(
	id valueObject.ResourceProfileId,
	name valueObject.ResourceProfileName,
	baseSpecs valueObject.ContainerSpecs,
	maxSpecs *valueObject.ContainerSpecs,
	scalingPolicy *valueObject.ScalingPolicy,
	scalingThreshold *uint64,
) ResourceProfile {
	return ResourceProfile{
		Id:               id,
		Name:             name,
		BaseSpecs:        baseSpecs,
		MaxSpecs:         maxSpecs,
		ScalingPolicy:    scalingPolicy,
		ScalingThreshold: scalingThreshold,
	}
}
