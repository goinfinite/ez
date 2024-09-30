package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/repository"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

func ReadRegistryTaggedImage(
	containerRegistryQueryRepo repository.ContainerRegistryQueryRepo,
	imageAddress valueObject.ContainerImageAddress,
) (taggedImage entity.RegistryTaggedImage, err error) {
	taggedImage, err = containerRegistryQueryRepo.ReadTaggedImage(imageAddress)
	if err != nil {
		slog.Error("ReadTaggedImageInfraError", slog.Any("error", err))
		return taggedImage, errors.New("ReadTaggedImageInfraError")
	}

	return taggedImage, nil
}
