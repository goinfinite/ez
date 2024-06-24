package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
)

func GetScheduledTasks(
	scheduledTaskQueryRepo repository.ScheduledTaskQueryRepo,
) ([]entity.ScheduledTask, error) {
	scheduledTasks, err := scheduledTaskQueryRepo.Get()
	if err != nil {
		log.Printf("GetTasksError: %s", err)
		return nil, errors.New("GetTasksInfraError")
	}

	return scheduledTasks, nil
}
