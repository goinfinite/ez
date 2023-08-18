package repository

import (
	"github.com/speedianet/sfm/src/domain/entity"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

type ContainerQueryRepo interface {
	Get() ([]entity.Container, error)
	GetById(id valueObject.ContainerId) (entity.Container, error)
}
