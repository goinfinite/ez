package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

func ImportContainerImageArchive(
	containerImageCmdRepo repository.ContainerImageCmdRepo,
	accountQueryRepo repository.AccountQueryRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	importDto dto.ImportContainerImageArchive,
) (imageId valueObject.ContainerImageId, err error) {
	accountEntity, err := accountQueryRepo.ReadFirst(dto.ReadAccountsRequest{
		AccountId: &importDto.AccountId,
	})
	if err != nil {
		slog.Error("ReadAccountInfoInfraError", slog.Any("error", err))
		return imageId, errors.New("ReadAccountInfoInfraError")
	}

	archiveFileSize, err := valueObject.NewByte(importDto.ArchiveFile.Size)
	if err != nil {
		return imageId, errors.New("InvalidArchiveFileSize")
	}

	accountStorageAvailable := accountEntity.Quota.StorageBytes - accountEntity.QuotaUsage.StorageBytes
	isThereStorageAvailable := accountStorageAvailable >= archiveFileSize
	if !isThereStorageAvailable {
		return imageId, errors.New("AccountStorageQuotaUsageExceeded")
	}

	imageId, err = containerImageCmdRepo.ImportArchive(importDto)
	if err != nil {
		slog.Error("ImportContainerImageArchiveInfraError", slog.Any("error", err))
		return imageId, errors.New("ImportContainerImageArchiveInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		ImportContainerImageArchive(importDto, imageId)

	return imageId, nil
}
