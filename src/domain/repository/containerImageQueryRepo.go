package repository

import (
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type ContainerImageQueryRepo interface {
	Read() ([]entity.ContainerImage, error)
	ReadById(
		accountId valueObject.AccountId,
		imageId valueObject.ContainerImageId,
	) (entity.ContainerImage, error)
}
