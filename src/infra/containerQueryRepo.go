package infra

import (
	"github.com/speedianet/sfm/src/domain/entity"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

type ContainerQueryRepo struct {
}

func (repo ContainerQueryRepo) getContainersByAccId(
	accId valueObject.AccountId,
) ([]entity.Container, error) {
	return []entity.Container{}, nil
}

func (repo ContainerQueryRepo) Get() ([]entity.Container, error) {
	allContainers := []entity.Container{}

	accsList, err := AccQueryRepo{}.Get()
	if err != nil {
		return allContainers, err
	}

	for _, acc := range accsList {
		accContainers, err := repo.getContainersByAccId(acc.Id)
		if err != nil {
			continue
		}
		allContainers = append(allContainers, accContainers...)
	}

	return allContainers, nil
}

func (repo ContainerQueryRepo) GetById(
	accId valueObject.AccountId,
	containerId valueObject.ContainerId,
) (entity.Container, error) {
	return entity.Container{}, nil
}
