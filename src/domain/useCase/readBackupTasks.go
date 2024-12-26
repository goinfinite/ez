package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
)

var BackupTasksDefaultPagination dto.Pagination = dto.Pagination{
	PageNumber:   0,
	ItemsPerPage: 10,
}

func ReadBackupTasks(
	backupQueryRepo repository.BackupQueryRepo,
	readDto dto.ReadBackupTasksRequest,
) (responseDto dto.ReadBackupTasksResponse, err error) {
	responseDto, err = backupQueryRepo.ReadTask(readDto)
	if err != nil {
		slog.Error("ReadBackupTasksInfraError", slog.Any("error", err))
		return responseDto, errors.New("ReadBackupTasksInfraError")
	}

	return responseDto, nil
}
