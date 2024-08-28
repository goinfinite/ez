package repository

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/valueObject"
)

type ContainerImageCmdRepo interface {
	CreateSnapshot(createDto dto.CreateContainerSnapshotImage) (
		valueObject.ContainerImageId, error,
	)
	Delete(deleteDto dto.DeleteContainerImage) error
}
