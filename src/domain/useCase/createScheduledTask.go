package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
)

func CreateScheduledTask(
	scheduledTaskCmdRepo repository.ScheduledTaskCmdRepo,
	dto dto.CreateScheduledTask,
) error {
	err := scheduledTaskCmdRepo.Create(dto)
	if err != nil {
		log.Printf("CreateScheduledTaskError: %v", err)
		return errors.New("CreateScheduledTaskInfraError")
	}

	return nil
}
