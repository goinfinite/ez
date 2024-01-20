package repository

import (
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type ContainerRegistryQueryRepo interface {
	GetImages(
		imageName *valueObject.RegistryImageName,
	) ([]entity.RegistryImage, error)
}
