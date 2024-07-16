package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
)

func DeleteActivityRecords(
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteActivityRecords,
) error {
	err := activityRecordCmdRepo.Delete(deleteDto)
	if err != nil {
		log.Printf("DeleteActivityRecordsError: %v", err)
		return errors.New("DeleteActivityRecordsInfraError")
	}

	return nil
}
