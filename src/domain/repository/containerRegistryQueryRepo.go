package repository

import (
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type ContainerRegistryQueryRepo interface {
	ReadImages(
		imageName *valueObject.RegistryImageName,
	) ([]entity.RegistryImage, error)
	ReadTaggedImage(
		imageAddress valueObject.ContainerImageAddress,
	) (entity.RegistryTaggedImage, error)
}
