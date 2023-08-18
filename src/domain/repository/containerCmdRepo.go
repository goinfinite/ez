package repository

import (
	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

type ContainerCmdRepo interface {
	Add(dto.AddContainer) error
	Delete(valueObject.ContainerId) error
}
