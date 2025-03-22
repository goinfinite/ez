package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
)

func UpdateBackupTask(
	backupQueryRepo repository.BackupQueryRepo,
	backupCmdRepo repository.BackupCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	updateDto dto.UpdateBackupTask,
) error {
	taskEntity, err := backupQueryRepo.ReadFirstTask(
		dto.ReadBackupTasksRequest{TaskId: &updateDto.TaskId},
	)
	if err != nil {
		return errors.New("BackupTaskNotFound")
	}

	err = backupCmdRepo.UpdateTask(updateDto)
	if err != nil {
		slog.Error("UpdateBackupTaskInfraError", slog.Any("error", err))
		return errors.New("UpdateBackupTaskInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		UpdateBackupTask(updateDto, taskEntity.AccountId)

	return nil
}
