package dto

import (
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type ContainerWithMetrics struct {
	entity.Container
	ResourceUsage valueObject.ContainerMetrics `json:"resourceUsage"`
}

func NewContainerWithMetrics(
	container entity.Container,
	containerMetrics valueObject.ContainerMetrics,
) ContainerWithMetrics {
	return ContainerWithMetrics{
		Container:     container,
		ResourceUsage: containerMetrics,
	}
}
