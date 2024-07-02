package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func ReadRegistryImages(
	containerRegistryQueryRepo repository.ContainerRegistryQueryRepo,
	imageName *valueObject.RegistryImageName,
) ([]entity.RegistryImage, error) {
	imagesList, err := containerRegistryQueryRepo.ReadImages(imageName)
	if err != nil {
		log.Printf("ReadImagesError: %s", err)
		return nil, errors.New("ReadImagesInfraError")
	}

	return imagesList, nil
}
