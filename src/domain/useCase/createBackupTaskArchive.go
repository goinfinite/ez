package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/repository"
)

var CreateBackupTaskArchiveDefaultTimeoutSecs uint32 = 21600

func CreateBackupTaskArchive(
	backupQueryRepo repository.BackupQueryRepo,
	backupCmdRepo repository.BackupCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateBackupTaskArchive,
) (taskArchiveEntity entity.BackupTaskArchive, err error) {
	taskEntity, err := backupQueryRepo.ReadFirstTask(
		dto.ReadBackupTasksRequest{TaskId: &createDto.TaskId},
	)
	if err != nil {
		slog.Error("BackupTaskNotFound", slog.String("error", err.Error()))
		return taskArchiveEntity, errors.New("BackupTaskNotFound")
	}

	taskArchiveId, err := backupCmdRepo.CreateTaskArchive(createDto)
	if err != nil {
		slog.Error("CreateBackupTaskArchiveError", slog.String("error", err.Error()))
		return taskArchiveEntity, errors.New("CreateBackupTaskArchiveInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateBackupTaskArchive(createDto, taskEntity.AccountId, taskArchiveId)

	requestDto := dto.ReadBackupTaskArchivesRequest{ArchiveId: &taskArchiveId}
	taskArchiveEntity, err = backupQueryRepo.ReadFirstTaskArchive(requestDto)
	if err != nil {
		slog.Error("ReadBackupTaskArchiveError", slog.String("error", err.Error()))
		return taskArchiveEntity, errors.New("ReadBackupTaskArchiveInfraError")
	}

	return taskArchiveEntity, nil
}
