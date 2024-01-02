package repository

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/valueObject"
)

type MappingCmdRepo interface {
	AddTargets(
		mappingId valueObject.MappingId,
		targets []valueObject.MappingTarget,
	) error
	Add(addMapping dto.AddMapping) error
	Delete(mappingId valueObject.MappingId) error
}
