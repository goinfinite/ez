package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

const BackupJobDefaultTimeoutSecs = valueObject.TimeDuration(24 * 60 * 60)

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
