package useCase

import (
	"errors"
	"log/slog"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
)

func ReadScheduledTasks(
	scheduledTaskQueryRepo repository.ScheduledTaskQueryRepo,
) ([]entity.ScheduledTask, error) {
	scheduledTasks, err := scheduledTaskQueryRepo.Read()
	if err != nil {
		slog.Error("GetTasksInfraError", slog.Any("error", err))
		return scheduledTasks, errors.New("GetTasksInfraError")
	}

	return scheduledTasks, nil
}
