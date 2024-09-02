package repository

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
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
	DeleteArchiveFile(dto.DeleteContainerImageArchiveFile) error
}
