package infra

import (
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/infra/db"
)

type ContainerImageQueryRepo struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewContainerImageQueryRepo(
	persistentDbSvc *db.PersistentDatabaseService,
) *ContainerImageQueryRepo {
	return &ContainerImageQueryRepo{persistentDbSvc: persistentDbSvc}
}

func (repo *ContainerImageQueryRepo) Read() ([]entity.ContainerImage, error) {
	return nil, nil
}
