package repository

import (
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type ContainerProfileQueryRepo interface {
	Get() ([]entity.ContainerProfile, error)
	GetById(id valueObject.ContainerProfileId) (entity.ContainerProfile, error)
}
