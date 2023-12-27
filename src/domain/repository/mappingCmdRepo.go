package repository

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/valueObject"
)

type MappingCmdRepo interface {
	Add(addMapping dto.AddMapping) error
	Delete(mappingId valueObject.MappingId) error
}
