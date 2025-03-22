package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
)

func DeleteBackupDestination(
	backupQueryRepo repository.BackupQueryRepo,
	backupCmdRepo repository.BackupCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteBackupDestination,
) error {
	destinationsReadRequestDto := dto.ReadBackupDestinationsRequest{
		DestinationId: &deleteDto.DestinationId,
		AccountId:     &deleteDto.AccountId,
	}
	_, err := backupQueryRepo.ReadFirstDestination(destinationsReadRequestDto, false)
	if err != nil {
		return errors.New("BackupDestinationNotFound")
	}

	jobsReadRequestDto := dto.ReadBackupJobsRequest{
		DestinationId: &deleteDto.DestinationId,
	}
	_, err = backupQueryRepo.ReadFirstJob(jobsReadRequestDto)
	if err == nil {
		return errors.New("CannotDeleteBackupDestinationInUse")
	}

	err = backupCmdRepo.DeleteDestination(deleteDto)
	if err != nil {
		slog.Error("DeleteBackupDestinationInfraError", slog.Any("error", err))
		return errors.New("DeleteBackupDestinationError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		DeleteBackupDestination(deleteDto)

	return nil
}
