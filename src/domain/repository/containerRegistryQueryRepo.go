package repository

import (
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type ContainerRegistryQueryRepo interface {
	ReadImages(
		imageName *valueObject.RegistryImageName,
	) ([]entity.RegistryImage, error)
	ReadTaggedImage(
		imageAddress valueObject.ContainerImageAddress,
	) (entity.RegistryTaggedImage, error)
}
