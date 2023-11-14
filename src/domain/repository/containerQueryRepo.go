package repository

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type ContainerQueryRepo interface {
	Get() ([]entity.Container, error)
	GetWithUsage() ([]dto.ContainerWithUsage, error)
	GetById(
		accId valueObject.AccountId,
		containerId valueObject.ContainerId,
	) (entity.Container, error)
	GetByAccId(accId valueObject.AccountId) ([]entity.Container, error)
}
