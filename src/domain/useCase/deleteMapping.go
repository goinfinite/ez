package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
)

func DeleteMapping(
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteMapping,
) error {
	_, err := mappingQueryRepo.ReadById(deleteDto.MappingId)
	if err != nil {
		return errors.New("MappingNotFound")
	}

	err = mappingCmdRepo.Delete(deleteDto)
	if err != nil {
		slog.Error("DeleteMappingInfraError", slog.Any("error", err))
		return errors.New("DeleteMappingError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		DeleteMapping(deleteDto)

	return nil
}
