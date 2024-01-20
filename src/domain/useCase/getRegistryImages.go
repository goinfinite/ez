package useCase

import (
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func GetRegistryImages(
	containerRegistryQueryRepo repository.ContainerRegistryQueryRepo,
	imageName *valueObject.RegistryImageName,
) ([]entity.RegistryImage, error) {
	return containerRegistryQueryRepo.GetImages(imageName)
}
