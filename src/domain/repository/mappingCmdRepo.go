package repository

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/valueObject"
)

type MappingCmdRepo interface {
	Add(addDto dto.AddMapping) (valueObject.MappingId, error)
	AddTarget(addDto dto.AddMappingTarget) error
	Delete(mappingId valueObject.MappingId) error
	DeleteTarget(targetId valueObject.MappingTargetId) error
}
