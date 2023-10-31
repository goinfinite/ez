package dto

import (
	"github.com/goinfinite/fleet/src/domain/entity"
	"github.com/goinfinite/fleet/src/domain/valueObject"
)

type ContainerWithUsage struct {
	entity.Container
	valueObject.ContainerResourceUsage
}
