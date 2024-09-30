package repository

import (
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type MappingCmdRepo interface {
	Create(dto.CreateMapping) (valueObject.MappingId, error)
	CreateTarget(dto.CreateMappingTarget) (valueObject.MappingTargetId, error)
	Delete(dto.DeleteMapping) error
	DeleteEmpty() error
	DeleteTarget(dto.DeleteMappingTarget) error
	DeleteTargetsByContainerId(valueObject.ContainerId) error
}
