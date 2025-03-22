package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
)

var BackupJobsDefaultPagination dto.Pagination = dto.Pagination{
	PageNumber:   0,
	ItemsPerPage: 10,
}

func ReadBackupJobs(
	backupQueryRepo repository.BackupQueryRepo,
	readDto dto.ReadBackupJobsRequest,
) (responseDto dto.ReadBackupJobsResponse, err error) {
	responseDto, err = backupQueryRepo.ReadJob(readDto)
	if err != nil {
		slog.Error("ReadBackupJobsInfraError", slog.Any("error", err))
		return responseDto, errors.New("ReadBackupJobsInfraError")
	}

	return responseDto, nil
}
