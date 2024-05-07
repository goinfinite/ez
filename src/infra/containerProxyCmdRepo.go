package infra

import (
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
)

type ContainerProxyCmdRepo struct {
	persistentDbSvc    *db.PersistentDatabaseService
	containerQueryRepo *ContainerQueryRepo
}

func NewContainerProxyCmdRepo(
	persistentDbSvc *db.PersistentDatabaseService,
) *ContainerProxyCmdRepo {
	return &ContainerProxyCmdRepo{
		persistentDbSvc:    persistentDbSvc,
		containerQueryRepo: NewContainerQueryRepo(persistentDbSvc),
	}
}

func (repo *ContainerProxyCmdRepo) Create(containerId valueObject.ContainerId) error {
	return nil
}

func (repo *ContainerProxyCmdRepo) Delete(containerId valueObject.ContainerId) error {
	return nil
}
