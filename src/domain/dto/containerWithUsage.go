package dto

import (
	"github.com/goinfinite/fleet/src/domain/entity"
	"github.com/goinfinite/fleet/src/domain/valueObject"
)

type ContainerWithUsage struct {
	entity.Container
	ResourceUsage valueObject.ContainerResourceUsage `json:"resourceUsage"`
}

func NewContainerWithUsage(
	container entity.Container,
	containerResourceUsage valueObject.ContainerResourceUsage,
) ContainerWithUsage {
	return ContainerWithUsage{
		Container:     container,
		ResourceUsage: containerResourceUsage,
	}
}
