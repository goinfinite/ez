package useCase

import (
	"log/slog"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

const ScheduledTasksRunIntervalSecs uint = 90

var scheduledTasksDefaultTimeoutSecs uint = 300

func RunScheduledTasks(
	scheduledTaskQueryRepo repository.ScheduledTaskQueryRepo,
	scheduledTaskCmdRepo repository.ScheduledTaskCmdRepo,
) {
	pendingStatus, _ := valueObject.NewScheduledTaskStatus("pending")
	readDto := dto.ReadScheduledTasksRequest{
		Pagination: ScheduledTasksDefaultPagination,
		TaskStatus: &pendingStatus,
	}

	responseDto, err := scheduledTaskQueryRepo.Read(readDto)
	if err != nil {
		slog.Error("ReadPendingScheduledTasksError", slog.Any("error", err))
		return
	}

	if len(responseDto.Tasks) == 0 {
		return
	}

	for _, pendingTask := range responseDto.Tasks {
		if pendingTask.RunAt != nil {
			nowUnixTime := valueObject.NewUnixTimeNow()
			if nowUnixTime.Read() < pendingTask.RunAt.Read() {
				continue
			}
		}

		if pendingTask.TimeoutSecs == nil {
			pendingTask.TimeoutSecs = &scheduledTasksDefaultTimeoutSecs
		}

		err = scheduledTaskCmdRepo.Run(pendingTask)
		if err != nil {
			slog.Error(
				"RunScheduledTaskError",
				slog.Uint64("scheduledTaskId", pendingTask.Id.Uint64()),
				slog.Any("error", err),
			)
			continue
		}
	}
}
