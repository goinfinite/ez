package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/repository"
)

func ReadContainerImages(
	containerImageQueryRepo repository.ContainerImageQueryRepo,
) ([]entity.ContainerImage, error) {
	containerImages, err := containerImageQueryRepo.Read()
	if err != nil {
		slog.Error("ReadContainerImagesInfraError", slog.Any("error", err))
		return containerImages, errors.New("ReadContainerImagesInfraError")
	}

	return containerImages, nil
}
