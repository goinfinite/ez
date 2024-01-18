package repository

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type ContainerQueryRepo interface {
	Get() ([]entity.Container, error)
	GetById(containerId valueObject.ContainerId) (entity.Container, error)
	GetByHostname(hostname valueObject.Fqdn) (entity.Container, error)
	GetByAccId(accId valueObject.AccountId) ([]entity.Container, error)
	GetWithMetrics() ([]dto.ContainerWithMetrics, error)
}
