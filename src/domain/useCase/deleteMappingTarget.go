package useCase

import (
	"errors"
	"log/slog"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
)

func DeleteMappingTarget(
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteMappingTarget,
) error {
	_, err := mappingQueryRepo.ReadTargetById(deleteDto.TargetId)
	if err != nil {
		return errors.New("MappingTargetNotFound")
	}

	err = mappingCmdRepo.DeleteTarget(deleteDto)
	if err != nil {
		slog.Error("DeleteMappingTargetInfraError", slog.Any("error", err))
		return errors.New("DeleteMappingTargetInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		DeleteMappingTarget(deleteDto)

	return nil
}
