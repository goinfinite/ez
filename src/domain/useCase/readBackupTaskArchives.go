package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
)

var BackupTaskArchivesDefaultPagination dto.Pagination = dto.Pagination{
	PageNumber:   0,
	ItemsPerPage: 10,
}

func ReadBackupTaskArchives(
	backupQueryRepo repository.BackupQueryRepo,
	requestDto dto.ReadBackupTaskArchivesRequest,
) (responseDto dto.ReadBackupTaskArchivesResponse, err error) {
	responseDto, err = backupQueryRepo.ReadTaskArchive(requestDto)
	if err != nil {
		slog.Error("ReadBackupTaskArchivesInfraError", slog.Any("error", err))
		return responseDto, errors.New("ReadBackupTaskArchivesInfraError")
	}

	return responseDto, nil
}
