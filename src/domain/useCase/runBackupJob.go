package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
)

func RunBackupJob(
	backupQueryRepo repository.BackupQueryRepo,
	backupCmdRepo repository.BackupCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	runDto dto.RunBackupJob,
) error {
	requestDto := dto.ReadBackupJobsRequest{
		JobId:     &runDto.JobId,
		AccountId: &runDto.AccountId,
	}

	_, err := backupQueryRepo.ReadFirstJob(requestDto)
	if err != nil {
		return errors.New("BackupJobNotFound")
	}

	err = backupCmdRepo.RunJob(runDto)
	if err != nil {
		slog.Error("RunBackupJobInfraError", slog.Any("error", err))
		return errors.New("RunBackupJobInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).RunBackupJob(runDto)

	return nil
}
