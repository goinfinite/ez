package useCase

import (
	"errors"
	"log/slog"

	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func DeleteMappingTarget(
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	mappingId valueObject.MappingId,
	targetId valueObject.MappingTargetId,
) error {
	_, err := mappingQueryRepo.ReadTargetById(targetId)
	if err != nil {
		return errors.New("MappingTargetNotFound")
	}

	err = mappingCmdRepo.DeleteTarget(mappingId, targetId)
	if err != nil {
		slog.Error("DeleteMappingTargetInfraError", slog.Any("error", err))
		return errors.New("DeleteMappingTargetInfraError")
	}

	slog.Info(
		"MappingTargetDeleted",
		slog.Uint64("mappingId", mappingId.Read()),
		slog.Uint64("targetId", targetId.Read()),
	)

	return nil
}
