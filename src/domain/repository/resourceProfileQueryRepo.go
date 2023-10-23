package repository

import (
	"github.com/speedianet/sfm/src/domain/entity"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

type ContainerProfileQueryRepo interface {
	Get() ([]entity.ContainerProfile, error)
	GetById(id valueObject.ContainerProfileId) (entity.ContainerProfile, error)
}
