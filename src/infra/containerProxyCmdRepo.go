package infra

import (
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
	dbModel "github.com/speedianet/control/src/infra/db/model"
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

func (repo *ContainerProxyCmdRepo) updateWebServerFile() error {
	return nil
}

func (repo *ContainerProxyCmdRepo) Create(containerId valueObject.ContainerId) error {
	containerEntity, err := repo.containerQueryRepo.GetById(containerId)
	if err != nil {
		return err
	}

	proxyModel := dbModel.NewContainerProxy(
		0,
		containerEntity.Id.String(),
		containerEntity.Hostname.String(),
	)

	createResult := repo.persistentDbSvc.Handler.Create(&proxyModel)
	if createResult.Error != nil {
		return createResult.Error
	}

	return repo.updateWebServerFile()
}

func (repo *ContainerProxyCmdRepo) Delete(containerId valueObject.ContainerId) error {
	err := repo.persistentDbSvc.Handler.Delete(
		dbModel.ContainerProxy{},
		"container_id = ?", containerId.String(),
	).Error
	if err != nil {
		return err
	}

	return repo.updateWebServerFile()
}
