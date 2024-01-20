package useCase

import (
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
)

func GetRegistryImages(
	containerRegistryQueryRepo repository.ContainerRegistryQueryRepo,
) ([]entity.RegistryImage, error) {
	return containerRegistryQueryRepo.GetImages()
}
