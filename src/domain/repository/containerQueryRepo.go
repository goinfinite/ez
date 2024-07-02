package repository

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type ContainerQueryRepo interface {
	Read() ([]entity.Container, error)
	ReadById(containerId valueObject.ContainerId) (entity.Container, error)
	ReadByHostname(hostname valueObject.Fqdn) (entity.Container, error)
	ReadByAccountId(accId valueObject.AccountId) ([]entity.Container, error)
	ReadWithMetrics() ([]dto.ContainerWithMetrics, error)
}
