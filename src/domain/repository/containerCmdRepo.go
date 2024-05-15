package repository

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/valueObject"
)

type ContainerCmdRepo interface {
	Create(createDto dto.CreateContainer) (valueObject.ContainerId, error)
	Update(updateDto dto.UpdateContainer) error
	Delete(
		accId valueObject.AccountId,
		containerId valueObject.ContainerId,
	) error
	GenerateContainerSessionToken(
		autoLoginDto dto.ContainerAutoLogin,
	) (valueObject.AccessTokenValue, error)
}
