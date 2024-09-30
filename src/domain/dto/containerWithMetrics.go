package dto

import (
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type ContainerWithMetrics struct {
	entity.Container
	Metrics valueObject.ContainerMetrics `json:"metrics"`
}

func NewContainerWithMetrics(
	container entity.Container,
	containerMetrics valueObject.ContainerMetrics,
) ContainerWithMetrics {
	return ContainerWithMetrics{
		Container: container,
		Metrics:   containerMetrics,
	}
}
