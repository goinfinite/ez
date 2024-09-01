package useCase

import (
	"errors"
	"log/slog"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
)

func CreateContainerImageArchiveFile(
	containerImageQueryRepo repository.ContainerImageQueryRepo,
	containerImageCmdRepo repository.ContainerImageCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateContainerImageArchiveFile,
) (archiveFile entity.ContainerImageArchiveFile, err error) {
	_, err = containerImageQueryRepo.ReadById(createDto.AccountId, createDto.ImageId)
	if err != nil {
		return archiveFile, errors.New("ContainerImageNotFound")
	}

	archiveFile, err = containerImageCmdRepo.CreateArchiveFile(createDto)
	if err != nil {
		slog.Error("CreateContainerImageArchiveFileInfraError", slog.Any("error", err))
		return archiveFile, errors.New("CreateContainerImageArchiveFileInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateContainerImageArchiveFile(createDto)

	return archiveFile, nil
}
