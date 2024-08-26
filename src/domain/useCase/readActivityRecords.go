package useCase

import (
	"log/slog"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
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
