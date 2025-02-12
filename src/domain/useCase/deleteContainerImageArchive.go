package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
)

func DeleteContainerImageArchive(
	containerImageQueryRepo repository.ContainerImageQueryRepo,
	containerImageCmdRepo repository.ContainerImageCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteContainerImageArchive,
) error {
	readDto := dto.NewReadContainerImageArchive(
		deleteDto.AccountId, deleteDto.ImageId,
	)
	archiveFile, err := containerImageQueryRepo.ReadArchive(readDto)
	if err != nil {
		return errors.New("ContainerImageArchiveNotFound")
	}

	err = containerImageCmdRepo.DeleteArchive(archiveFile)
	if err != nil {
		slog.Error("DeleteContainerImageArchiveInfraError", slog.Any("error", err))
		return errors.New("DeleteContainerImageArchiveError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		DeleteContainerImageArchive(deleteDto)

	return nil
}
