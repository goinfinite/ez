package repository

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/valueObject"
)

type ContainerProfileCmdRepo interface {
	Create(dto.CreateContainerProfile) (valueObject.ContainerProfileId, error)
	Update(dto.UpdateContainerProfile) error
	Delete(dto.DeleteContainerProfile) error
}
