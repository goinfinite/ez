package repository

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/valueObject"
)

type MappingCmdRepo interface {
	Create(createDto dto.CreateMapping) (valueObject.MappingId, error)
	CreateTarget(createDto dto.CreateMappingTarget) error
	Delete(mappingId valueObject.MappingId) error
	DeleteTarget(
		mappingId valueObject.MappingId,
		targetId valueObject.MappingTargetId,
	) error

	CreateContainerProxy(containerId valueObject.ContainerId) error
}
