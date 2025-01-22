package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
)

func DeleteBackupTask(
	backupQueryRepo repository.BackupQueryRepo,
	backupCmdRepo repository.BackupCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteBackupTask,
) error {
	tasksReadRequestDto := dto.ReadBackupTasksRequest{
		TaskId: &deleteDto.TaskId,
	}
	taskEntity, err := backupQueryRepo.ReadFirstTask(tasksReadRequestDto)
	if err != nil {
		return errors.New("BackupTaskNotFound")
	}

	err = backupCmdRepo.DeleteTask(deleteDto)
	if err != nil {
		slog.Error("DeleteBackupTaskInfraError", slog.Any("error", err))
		return errors.New("DeleteBackupTaskError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		DeleteBackupTask(deleteDto, taskEntity.AccountId)

	return nil
}
