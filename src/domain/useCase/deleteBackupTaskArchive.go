package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
)

func DeleteBackupTaskArchive(
	backupQueryRepo repository.BackupQueryRepo,
	backupCmdRepo repository.BackupCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteBackupTaskArchive,
) error {
	taskArchiveEntity, err := backupQueryRepo.ReadFirstTaskArchive(
		dto.ReadBackupTaskArchivesRequest{ArchiveId: &deleteDto.ArchiveId},
	)
	if err != nil {
		return errors.New("BackupTaskArchiveNotFound")
	}

	err = backupCmdRepo.DeleteTaskArchive(deleteDto)
	if err != nil {
		slog.Error("DeleteBackupTaskArchiveError", slog.String("error", err.Error()))
		return errors.New("DeleteBackupTaskArchiveInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		DeleteBackupTaskArchive(deleteDto, taskArchiveEntity.AccountId)

	return nil
}
