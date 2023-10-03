package infra

import (
	"github.com/speedianet/sfm/src/domain/entity"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

type ResourceProfileQueryRepo struct {
}

func (repo ResourceProfileQueryRepo) Get() ([]entity.ResourceProfile, error) {
	return []entity.ResourceProfile{}, nil
}

func (repo ResourceProfileQueryRepo) GetById(
	id valueObject.ResourceProfileId,
) (entity.ResourceProfile, error) {
	return entity.ResourceProfile{}, nil
}
