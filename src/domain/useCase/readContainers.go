package useCase

import (
	"errors"
	"log/slog"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
)

func ReadContainers(
	containerQueryRepo repository.ContainerQueryRepo,
) ([]entity.Container, error) {
	containersList, err := containerQueryRepo.Read()
	if err != nil {
		slog.Error("ReadContainersInfraError", slog.Any("error", err))
		return containersList, errors.New("ReadContainersInfraError")
	}

	return containersList, nil
}
