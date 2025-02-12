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

	CreateArchive(
		dto.CreateContainerImageArchive,
	) (entity.ContainerImageArchive, error)
	ImportArchive(dto.ImportContainerImageArchive) (
		valueObject.ContainerImageId, error,
	)
	DeleteArchive(entity.ContainerImageArchive) error
}
