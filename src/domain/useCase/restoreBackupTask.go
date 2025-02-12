package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
)

var RestoreBackupTaskDefaultTimeoutSecs uint32 = 21600

func RestoreBackupTask(
	backupQueryRepo repository.BackupQueryRepo,
	backupCmdRepo repository.BackupCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	restoreDto dto.RestoreBackupTask,
) error {
	taskEntity, err := backupQueryRepo.ReadFirstTask(
		dto.ReadBackupTasksRequest{TaskId: &restoreDto.TaskId},
	)
	if err != nil {
		slog.Error("BackupTaskNotFound", slog.String("error", err.Error()))
		return errors.New("BackupTaskNotFound")
	}

	err = backupCmdRepo.RestoreTask(restoreDto)
	if err != nil {
		slog.Error("RestoreTaskError", slog.String("error", err.Error()))
		return errors.New("RestoreTaskInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		RestoreBackupTask(restoreDto, taskEntity)

	return nil
}
