package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
)

var CreateBackupTaskArchiveDefaultTimeoutSecs uint32 = 21600

func CreateBackupTaskArchive(
	backupQueryRepo repository.BackupQueryRepo,
	backupCmdRepo repository.BackupCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateBackupTaskArchive,
) error {
	taskEntity, err := backupQueryRepo.ReadFirstTask(
		dto.ReadBackupTasksRequest{TaskId: &createDto.TaskId},
	)
	if err != nil {
		slog.Error("BackupTaskNotFound", slog.String("error", err.Error()))
		return errors.New("BackupTaskNotFound")
	}

	taskArchiveId, err := backupCmdRepo.CreateTaskArchive(createDto)
	if err != nil {
		slog.Error("CreateBackupTaskArchiveInfraError", slog.String("error", err.Error()))
		return errors.New("CreateBackupTaskArchiveInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateBackupTaskArchive(createDto, taskEntity.AccountId, taskArchiveId)

	return nil
}
