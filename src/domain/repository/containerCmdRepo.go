package repository

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/valueObject"
)

type ContainerCmdRepo interface {
	Add(addDto dto.AddContainer) error
	Update(updateDto dto.UpdateContainer) error
	Delete(
		accId valueObject.AccountId,
		containerId valueObject.ContainerId,
	) error
}
