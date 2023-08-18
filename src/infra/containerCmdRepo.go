package infra

import (
	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

type ContainerCmdRepo struct {
}

func (repo ContainerCmdRepo) Add(dto.AddContainer) error {
	return nil
}

func (repo ContainerCmdRepo) Update(dto.UpdateContainer) error {
	return nil
}

func (repo ContainerCmdRepo) Delete(valueObject.ContainerId) error {
	return nil
}
