package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/repository"
)

func ReadContainerProfiles(
	containerProfileQueryRepo repository.ContainerProfileQueryRepo,
) ([]entity.ContainerProfile, error) {
	profilesList, err := containerProfileQueryRepo.Read()
	if err != nil {
		slog.Error("ReadContainerProfilesInfraError", slog.Any("error", err))
		return profilesList, errors.New("ReadContainerProfilesInfraError")
	}

	return profilesList, nil
}
