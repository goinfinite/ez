package useCase

import (
	"errors"
	"log/slog"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
)

func ReadContainerImageArchiveFiles(
	containerImageQueryRepo repository.ContainerImageQueryRepo,
) ([]entity.ContainerImageArchiveFile, error) {
	archiveFiles, err := containerImageQueryRepo.ReadArchiveFiles()
	if err != nil {
		slog.Error("ReadContainerImageArchiveFilesInfraError", slog.Any("error", err))
		return archiveFiles, errors.New("ReadContainerImageArchiveFilesInfraError")
	}

	return archiveFiles, nil
}
