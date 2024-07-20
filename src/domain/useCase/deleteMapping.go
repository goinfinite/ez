package useCase

import (
	"errors"
	"log/slog"

	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func DeleteMapping(
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	mappingId valueObject.MappingId,
) error {
	_, err := mappingQueryRepo.ReadById(mappingId)
	if err != nil {
		return errors.New("MappingNotFound")
	}

	err = mappingCmdRepo.Delete(mappingId)
	if err != nil {
		slog.Error("DeleteMappingInfraError", slog.Any("error", err))
		return errors.New("DeleteMappingError")
	}

	slog.Info("MappingDeleted", slog.Uint64("mappingId", mappingId.Uint64()))
	return nil
}
