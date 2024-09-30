package useCase

import (
	"errors"
	"log"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
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
