package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
)

func UpdateBackupDestination(
	backupQueryRepo repository.BackupQueryRepo,
	backupCmdRepo repository.BackupCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	updateDto dto.UpdateBackupDestination,
) error {
	requestDto := dto.ReadBackupDestinationsRequest{
		DestinationId: &updateDto.DestinationId,
		AccountId:     &updateDto.AccountId,
	}

	_, err := backupQueryRepo.ReadFirstDestination(requestDto, false)
	if err != nil {
		return errors.New("BackupDestinationNotFound")
	}

	err = backupCmdRepo.UpdateDestination(updateDto)
	if err != nil {
		slog.Error("UpdateBackupDestinationInfraError", slog.Any("error", err))
		return errors.New("UpdateBackupDestinationInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		UpdateBackupDestination(updateDto)

	return nil
}
