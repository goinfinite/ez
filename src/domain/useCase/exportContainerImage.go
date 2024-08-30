package useCase

import (
	"errors"
	"log/slog"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func ExportContainerImage(
	containerImageQueryRepo repository.ContainerImageQueryRepo,
	containerImageCmdRepo repository.ContainerImageCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	exportDto dto.ExportContainerImage,
) (downloadUrl valueObject.Url, err error) {
	_, err = containerImageQueryRepo.ReadById(exportDto.AccountId, exportDto.ImageId)
	if err != nil {
		return downloadUrl, errors.New("ContainerImageNotFound")
	}

	downloadUrl, err = containerImageCmdRepo.Export(exportDto)
	if err != nil {
		slog.Error("ExportContainerImageInfraError", slog.Any("error", err))
		return downloadUrl, errors.New("ExportContainerImageError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).ExportContainerImage(exportDto)

	return downloadUrl, nil
}
