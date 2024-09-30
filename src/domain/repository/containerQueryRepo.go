package repository

import (
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type ContainerQueryRepo interface {
	Read() ([]entity.Container, error)
	ReadById(valueObject.ContainerId) (entity.Container, error)
	ReadByHostname(valueObject.Fqdn) (entity.Container, error)
	ReadByAccountId(valueObject.AccountId) ([]entity.Container, error)
	ReadByImageId(valueObject.AccountId, valueObject.ContainerImageId) ([]entity.Container, error)
	ReadWithMetrics() ([]dto.ContainerWithMetrics, error)
	ReadWithMetricsById(valueObject.ContainerId) (dto.ContainerWithMetrics, error)
}
