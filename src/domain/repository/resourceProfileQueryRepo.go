package repository

import (
	"github.com/speedianet/sfm/src/domain/entity"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

type ResourceProfileQueryRepo interface {
	Get() ([]entity.ResourceProfile, error)
	GetById(id valueObject.ResourceProfileId) (entity.ResourceProfile, error)
}
