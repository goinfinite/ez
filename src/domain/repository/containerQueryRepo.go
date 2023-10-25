package repository

import (
	"github.com/goinfinite/fleet/src/domain/entity"
	"github.com/goinfinite/fleet/src/domain/valueObject"
)

type ContainerQueryRepo interface {
	Get() ([]entity.Container, error)
	GetById(
		accId valueObject.AccountId,
		containerId valueObject.ContainerId,
	) (entity.Container, error)
}
