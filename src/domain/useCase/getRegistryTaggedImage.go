package useCase

import (
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func GetRegistryTaggedImage(
	containerRegistryQueryRepo repository.ContainerRegistryQueryRepo,
	imageAddress valueObject.ContainerImageAddress,
) (entity.RegistryTaggedImage, error) {
	return containerRegistryQueryRepo.GetTaggedImage(imageAddress)
}
