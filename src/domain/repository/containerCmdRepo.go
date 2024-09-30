package repository

import (
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type ContainerCmdRepo interface {
	Create(dto.CreateContainer) (valueObject.ContainerId, error)
	Update(dto.UpdateContainer) error
	Delete(dto.DeleteContainer) error
	CreateContainerSessionToken(
		dto.CreateContainerSessionToken,
	) (valueObject.AccessTokenValue, error)
}
