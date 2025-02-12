package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/repository"
)

func ReadContainerImageArchive(
	containerImageQueryRepo repository.ContainerImageQueryRepo,
	readDto dto.ReadContainerImageArchive,
) (archiveFile entity.ContainerImageArchive, err error) {
	archiveFile, err = containerImageQueryRepo.ReadArchive(readDto)
	if err != nil {
		slog.Error("ReadContainerImageArchiveInfraError", slog.Any("error", err))
		return archiveFile, errors.New("ReadContainerImageArchiveInfraError")
	}

	return archiveFile, nil
}
