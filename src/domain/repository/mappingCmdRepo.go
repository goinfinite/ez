package repository

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/valueObject"
)

type MappingCmdRepo interface {
	Create(dto.CreateMapping) (valueObject.MappingId, error)
	CreateTarget(dto.CreateMappingTarget) (valueObject.MappingTargetId, error)
	Delete(dto.DeleteMapping) error
	DeleteEmpty() error
	DeleteTarget(dto.DeleteMappingTarget) error
	DeleteTargetsByContainerId(valueObject.ContainerId) error
}
