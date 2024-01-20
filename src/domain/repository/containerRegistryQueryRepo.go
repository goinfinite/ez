package repository

import (
	"github.com/speedianet/control/src/domain/entity"
)

type ContainerRegistryQueryRepo interface {
	GetImages() ([]entity.RegistryImage, error)
}
