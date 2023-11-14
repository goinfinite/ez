package dto

import (
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
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
