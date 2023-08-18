package repository

import (
	"github.com/speedianet/sfm/src/domain/dto"
)

type ContainerCmdRepo interface {
	Add(dto.AddContainer) error
}
