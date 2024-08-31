package useCase

import (
	"errors"
	"log/slog"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
)

func ExportContainerImage(
	containerImageQueryRepo repository.ContainerImageQueryRepo,
	containerImageCmdRepo repository.ContainerImageCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	exportDto dto.ExportContainerImage,
) (archiveFile entity.ContainerImageArchiveFile, err error) {
	_, err = containerImageQueryRepo.ReadById(exportDto.AccountId, exportDto.ImageId)
	if err != nil {
		return archiveFile, errors.New("ContainerImageNotFound")
	}

	archiveFile, err = containerImageCmdRepo.Export(exportDto)
	if err != nil {
		slog.Error("ExportContainerImageInfraError", slog.Any("error", err))
		return archiveFile, errors.New("ExportContainerImageError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).ExportContainerImage(exportDto)

	return archiveFile, nil
}
