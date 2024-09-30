package repository

import (
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type MappingQueryRepo interface {
	Read() ([]entity.Mapping, error)
	ReadById(id valueObject.MappingId) (entity.Mapping, error)
	ReadTargetById(id valueObject.MappingTargetId) (entity.MappingTarget, error)
	ReadTargetsByContainerId(
		containerId valueObject.ContainerId,
	) ([]entity.MappingTarget, error)
}
