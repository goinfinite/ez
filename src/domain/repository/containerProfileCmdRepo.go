package repository

import (
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type ContainerProfileCmdRepo interface {
	Create(dto.CreateContainerProfile) (valueObject.ContainerProfileId, error)
	Update(dto.UpdateContainerProfile) error
	Delete(dto.DeleteContainerProfile) error
}
