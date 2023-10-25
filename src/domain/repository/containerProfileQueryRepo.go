package repository

import (
	"github.com/goinfinite/fleet/src/domain/entity"
	"github.com/goinfinite/fleet/src/domain/valueObject"
)

type ContainerProfileQueryRepo interface {
	Get() ([]entity.ContainerProfile, error)
	GetById(id valueObject.ContainerProfileId) (entity.ContainerProfile, error)
}
