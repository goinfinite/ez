package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/repository"
)

func ReadBackupTaskArchive(
	backupQueryRepo repository.BackupQueryRepo,
	requestDto dto.ReadBackupTaskArchivesRequest,
) (archiveFile entity.BackupTaskArchive, err error) {
	archiveFile, err = backupQueryRepo.ReadFirstTaskArchive(requestDto)
	if err != nil {
		slog.Error("ReadBackupTaskArchiveFileError", slog.String("error", err.Error()))
		return archiveFile, errors.New("ReadBackupTaskArchiveFileInfraError")
	}

	return archiveFile, nil
}
