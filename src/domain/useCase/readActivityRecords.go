package useCase

import (
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/repository"
)

func ReadActivityRecords(
	activityRecordQueryRepo repository.ActivityRecordQueryRepo,
	readDto dto.ReadActivityRecords,
) (activityRecords []entity.ActivityRecord) {
	activityRecords, err := activityRecordQueryRepo.Read(readDto)
	if err != nil {
		slog.Error("ReadActivityRecordsInfraError", slog.Any("error", err))
	}

	return activityRecords
}
