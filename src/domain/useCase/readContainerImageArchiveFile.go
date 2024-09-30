package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/repository"
)

func ReadContainerImageArchiveFile(
	containerImageQueryRepo repository.ContainerImageQueryRepo,
	readDto dto.ReadContainerImageArchiveFile,
) (archiveFile entity.ContainerImageArchiveFile, err error) {
	archiveFile, err = containerImageQueryRepo.ReadArchiveFile(readDto)
	if err != nil {
		slog.Error("ReadContainerImageArchiveFileInfraError", slog.Any("error", err))
		return archiveFile, errors.New("ReadContainerImageArchiveFileInfraError")
	}

	return archiveFile, nil
}
