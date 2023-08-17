package dto

import "github.com/speedianet/sfm/src/domain/valueObject"

type UpdateContainer struct {
	ContainerId valueObject.ContainerId     `json:"id"`
	Status      bool                        `json:"status"`
	BaseSpecs   *valueObject.ContainerSpecs `json:"baseSpecs"`
	MaxSpecs    *valueObject.ContainerSpecs `json:"maxSpecs"`
}

func NewUpdateContainer(
	containerId valueObject.ContainerId,
	status bool,
	baseSpecs *valueObject.ContainerSpecs,
	maxSpecs *valueObject.ContainerSpecs,
) UpdateContainer {
	return UpdateContainer{
		ContainerId: containerId,
		Status:      status,
		BaseSpecs:   baseSpecs,
		MaxSpecs:    maxSpecs,
	}
}
