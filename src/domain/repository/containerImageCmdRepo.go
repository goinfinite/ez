package repository

import (
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type ContainerImageCmdRepo interface {
	CreateSnapshot(dto.CreateContainerSnapshotImage) (
		valueObject.ContainerImageId, error,
	)
	Delete(dto.DeleteContainerImage) error

	CreateArchiveFile(
		dto.CreateContainerImageArchiveFile,
	) (entity.ContainerImageArchiveFile, error)
	ImportArchiveFile(dto.ImportContainerImageArchiveFile) (
		valueObject.ContainerImageId, error,
	)
	DeleteArchiveFile(entity.ContainerImageArchiveFile) error
}
