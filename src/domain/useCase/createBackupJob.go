package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
)

func CreateBackupJob(
	backupCmdRepo repository.BackupCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateBackupJob,
) error {
	backupJobId, err := backupCmdRepo.CreateJob(createDto)
	if err != nil {
		slog.Error("CreateBackupJobInfraError", slog.Any("error", err))
		return errors.New("CreateBackupJobInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateBackupJob(createDto, backupJobId)

	return nil
}
