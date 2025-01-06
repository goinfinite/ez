package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
)

func DeleteBackupJob(
	backupQueryRepo repository.BackupQueryRepo,
	backupCmdRepo repository.BackupCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteBackupJob,
) error {
	jobsReadRequestDto := dto.ReadBackupJobsRequest{
		JobId:     &deleteDto.JobId,
		AccountId: &deleteDto.AccountId,
	}
	_, err := backupQueryRepo.ReadFirstJob(jobsReadRequestDto)
	if err != nil {
		return errors.New("BackupJobNotFound")
	}

	err = backupCmdRepo.DeleteJob(deleteDto)
	if err != nil {
		slog.Error("DeleteBackupJobInfraError", slog.Any("error", err))
		return errors.New("DeleteBackupJobError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		DeleteBackupJob(deleteDto)

	return nil
}
