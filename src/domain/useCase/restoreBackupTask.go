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
	requestRestoreDto dto.RestoreBackupTaskRequest,
) error {
	if requestRestoreDto.ArchiveId == nil && requestRestoreDto.TaskId == nil {
		return errors.New("TaskIdOrArchiveIdRequired")
	}

	accountId := requestRestoreDto.OperatorAccountId
	if requestRestoreDto.ArchiveId == nil {
		taskEntity, err := backupQueryRepo.ReadFirstTask(
			dto.ReadBackupTasksRequest{TaskId: requestRestoreDto.TaskId},
		)
		if err != nil {
			slog.Error("BackupTaskNotFound", slog.String("error", err.Error()))
			return errors.New("BackupTaskNotFound")
		}
		accountId = taskEntity.AccountId
	}

	err := backupCmdRepo.RestoreTask(requestRestoreDto)
	if err != nil {
		slog.Error("RestoreTaskError", slog.String("error", err.Error()))
		return errors.New("RestoreTaskInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		RestoreBackupTask(requestRestoreDto, accountId)

	return nil
}
