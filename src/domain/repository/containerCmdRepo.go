package repository

import (
	"github.com/goinfinite/fleet/src/domain/dto"
	"github.com/goinfinite/fleet/src/domain/entity"
	"github.com/goinfinite/fleet/src/domain/valueObject"
)

type ContainerCmdRepo interface {
	Add(addContainer dto.AddContainer) error
	Update(
		currentContainer entity.Container,
		updateContainer dto.UpdateContainer,
	) error
	Delete(
		accId valueObject.AccountId,
		containerId valueObject.ContainerId,
	) error
}
