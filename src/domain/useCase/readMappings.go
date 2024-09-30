package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/repository"
)

func ReadMappings(
	mappingQueryRepo repository.MappingQueryRepo,
) ([]entity.Mapping, error) {
	mappingsList, err := mappingQueryRepo.Read()
	if err != nil {
		slog.Error("ReadMappingsInfraError", slog.Any("error", err))
		return mappingsList, errors.New("ReadMappingsInfraError")
	}

	return mappingsList, nil
}
