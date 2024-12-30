package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
)

var BackupDestinationsDefaultPagination dto.Pagination = dto.Pagination{
	PageNumber:   0,
	ItemsPerPage: 10,
}

func ReadBackupDestinations(
	backupQueryRepo repository.BackupQueryRepo,
	readDto dto.ReadBackupDestinationsRequest,
) (responseDto dto.ReadBackupDestinationsResponse, err error) {
	responseDto, err = backupQueryRepo.ReadDestination(readDto, false)
	if err != nil {
		slog.Error("ReadBackupDestinationsInfraError", slog.Any("error", err))
		return responseDto, errors.New("ReadBackupDestinationsInfraError")
	}

	return responseDto, nil
}
