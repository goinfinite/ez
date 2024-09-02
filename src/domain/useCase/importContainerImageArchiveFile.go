package useCase

import (
	"errors"
	"log/slog"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func ImportContainerImageArchiveFile(
	containerImageCmdRepo repository.ContainerImageCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	importDto dto.ImportContainerImageArchiveFile,
) (imageId valueObject.ContainerImageId, err error) {
	imageId, err = containerImageCmdRepo.ImportArchiveFile(importDto)
	if err != nil {
		slog.Error("ImportContainerImageArchiveFileInfraError", slog.Any("error", err))
		return imageId, errors.New("ImportContainerImageArchiveFileInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		ImportContainerImageArchiveFile(importDto)

	return imageId, nil
}
