package infra

import (
	"strings"

	"github.com/speedianet/sfm/src/domain/entity"
	"github.com/speedianet/sfm/src/domain/valueObject"
	infraHelper "github.com/speedianet/sfm/src/infra/helper"
)

type ContainerQueryRepo struct {
}

func (repo ContainerQueryRepo) GetById(
	accId valueObject.AccountId,
	containerId valueObject.ContainerId,
) (entity.Container, error) {
	return entity.Container{}, nil
}

func (repo ContainerQueryRepo) GetByAccId(
	accId valueObject.AccountId,
) ([]entity.Container, error) {
	containersIds, err := infraHelper.RunCmdAsUser(
		accId,
		"podman",
		"container",
		"list",
		"--format '{{.ID}}'",
	)
	if err != nil {
		return []entity.Container{}, err
	}
	containersIdsList := strings.Split(containersIds, "\n")
	if len(containersIdsList) == 0 {
		return []entity.Container{}, nil
	}

	containers := []entity.Container{}
	for _, containerId := range containersIdsList {
		containerId = strings.TrimSpace(containerId)
		containerIdValidated, err := valueObject.NewContainerId(containerId)
		if err != nil {
			continue
		}

		container, err := repo.GetById(accId, containerIdValidated)
		if err != nil {
			continue
		}
		containers = append(containers, container)
	}

	return containers, nil
}

func (repo ContainerQueryRepo) Get() ([]entity.Container, error) {
	allContainers := []entity.Container{}

	accsList, err := AccQueryRepo{}.Get()
	if err != nil {
		return allContainers, err
	}

	for _, acc := range accsList {
		accContainers, err := repo.GetByAccId(acc.Id)
		if err != nil {
			continue
		}
		allContainers = append(allContainers, accContainers...)
	}

	return allContainers, nil
}
