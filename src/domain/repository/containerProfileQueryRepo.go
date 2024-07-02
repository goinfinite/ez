package repository

import (
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type ContainerProfileQueryRepo interface {
	Read() ([]entity.ContainerProfile, error)
	ReadById(id valueObject.ContainerProfileId) (entity.ContainerProfile, error)
}
