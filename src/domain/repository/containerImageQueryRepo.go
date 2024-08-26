package repository

import (
	"github.com/speedianet/control/src/domain/entity"
)

type ContainerImageQueryRepo interface {
	Read() ([]entity.ContainerImage, error)
}
