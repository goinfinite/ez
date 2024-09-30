package repository

import (
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type ContainerProfileQueryRepo interface {
	Read() ([]entity.ContainerProfile, error)
	ReadById(id valueObject.ContainerProfileId) (entity.ContainerProfile, error)
}
