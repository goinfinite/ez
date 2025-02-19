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
	if restoreDto.ArchiveId == nil && restoreDto.TaskId == nil {
		return errors.New("TaskIdOrArchiveIdRequired")
	}

	accountId := restoreDto.OperatorAccountId
	if restoreDto.ArchiveId == nil {
		taskEntity, err := backupQueryRepo.ReadFirstTask(
			dto.ReadBackupTasksRequest{TaskId: restoreDto.TaskId},
		)
		if err != nil {
			slog.Error("BackupTaskNotFound", slog.String("error", err.Error()))
			return errors.New("BackupTaskNotFound")
		}
		accountId = taskEntity.AccountId
	}

	err := backupCmdRepo.RestoreTask(restoreDto)
	if err != nil {
		slog.Error("RestoreTaskError", slog.String("error", err.Error()))
		return errors.New("RestoreTaskInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		RestoreBackupTask(restoreDto, accountId, restoreDto.TaskId, restoreDto.ArchiveId)

	return nil
}
