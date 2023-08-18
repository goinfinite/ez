package repository

import (
	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

type ContainerCmdRepo interface {
	Add(dto.AddContainer) error
	Update(dto.UpdateContainer) error
	Delete(valueObject.ContainerId) error
}
