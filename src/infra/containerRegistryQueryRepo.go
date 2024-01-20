package infra

import (
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/infra/db"
)

type ContainerRegistryQueryRepo struct {
	dbSvc *db.DatabaseService
}

func NewContainerRegistryQueryRepo(dbSvc *db.DatabaseService) *ContainerRegistryQueryRepo {
	return &ContainerRegistryQueryRepo{dbSvc: dbSvc}
}

func (repo ContainerRegistryQueryRepo) GetImages() ([]entity.RegistryImage, error) {
	return []entity.RegistryImage{}, nil
}
