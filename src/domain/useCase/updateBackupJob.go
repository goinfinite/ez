package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
)

func UpdateBackupJob(
	backupQueryRepo repository.BackupQueryRepo,
	backupCmdRepo repository.BackupCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	updateDto dto.UpdateBackupJob,
) error {
	requestDto := dto.ReadBackupJobsRequest{
		JobId:     &updateDto.JobId,
		AccountId: &updateDto.AccountId,
	}

	_, err := backupQueryRepo.ReadFirstJob(requestDto)
	if err != nil {
		return errors.New("BackupJobNotFound")
	}

	err = backupCmdRepo.UpdateJob(updateDto)
	if err != nil {
		slog.Error("UpdateBackupJobInfraError", slog.Any("error", err))
		return errors.New("UpdateBackupJobInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		UpdateBackupJob(updateDto)

	return nil
}
