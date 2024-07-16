package useCase

import (
	"log"

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
		log.Printf("ReadActivityRecordsInfraError: %s", err)
	}

	return activityRecords
}
