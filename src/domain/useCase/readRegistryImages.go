package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/repository"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

func ReadRegistryImages(
	containerRegistryQueryRepo repository.ContainerRegistryQueryRepo,
	imageName *valueObject.RegistryImageName,
) ([]entity.RegistryImage, error) {
	imagesList, err := containerRegistryQueryRepo.ReadImages(imageName)
	if err != nil {
		slog.Error("ReadImagesInfraError", slog.Any("error", err))
		return imagesList, errors.New("ReadImagesInfraError")
	}

	return imagesList, nil
}
