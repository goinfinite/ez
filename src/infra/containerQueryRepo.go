package infra

import (
	"github.com/speedianet/sfm/src/domain/entity"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

type ContainerQueryRepo struct {
}

func (repo ContainerQueryRepo) Get() ([]entity.Container, error) {
	return []entity.Container{}, nil
}

func (repo ContainerQueryRepo) GetById(
	accId valueObject.AccountId,
	containerId valueObject.ContainerId,
) (entity.Container, error) {
	return entity.Container{}, nil
}
