package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
)

var ContainersDefaultPagination dto.Pagination = dto.Pagination{
	PageNumber:   0,
	ItemsPerPage: 10,
}

func ReadContainers(
	containerQueryRepo repository.ContainerQueryRepo,
	readDto dto.ReadContainersRequest,
) (responseDto dto.ReadContainersResponse, err error) {
	responseDto, err = containerQueryRepo.Read(readDto)
	if err != nil {
		slog.Error("ReadContainersInfraError", slog.Any("error", err))
		return responseDto, errors.New("ReadContainersInfraError")
	}

	return responseDto, nil
}
