package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
)

var ContainerImageArchivesDefaultPagination dto.Pagination = dto.Pagination{
	PageNumber:   0,
	ItemsPerPage: 5,
}

func ReadContainerImageArchives(
	containerImageQueryRepo repository.ContainerImageQueryRepo,
	requestDto dto.ReadContainerImageArchivesRequest,
) (dto.ReadContainerImageArchivesResponse, error) {
	responseDto, err := containerImageQueryRepo.ReadArchives(requestDto)
	if err != nil {
		slog.Error("ReadContainerImageArchivesError", slog.String("error", err.Error()))
		return responseDto, errors.New("ReadContainerImageArchivesInfraError")
	}

	return responseDto, nil
}
