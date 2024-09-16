package useCase

import (
	"errors"
	"log/slog"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
)

func DeleteContainerImageArchiveFile(
	containerImageQueryRepo repository.ContainerImageQueryRepo,
	containerImageCmdRepo repository.ContainerImageCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteContainerImageArchiveFile,
) error {
	readDto := dto.NewReadContainerImageArchiveFile(
		deleteDto.AccountId, deleteDto.ImageId,
	)
	archiveFile, err := containerImageQueryRepo.ReadArchiveFile(readDto)
	if err != nil {
		return errors.New("ContainerImageArchiveFileNotFound")
	}

	err = containerImageCmdRepo.DeleteArchiveFile(archiveFile)
	if err != nil {
		slog.Error("DeleteContainerImageArchiveFileInfraError", slog.Any("error", err))
		return errors.New("DeleteContainerImageArchiveFileError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		DeleteContainerImageArchiveFile(deleteDto)

	return nil
}
