package repository

import (
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type MappingQueryRepo interface {
	Get() ([]entity.Mapping, error)
	GetById(id valueObject.MappingId) (entity.Mapping, error)
	GetByHostname(hostname valueObject.Fqdn) ([]entity.Mapping, error)
	GetTargetById(id valueObject.MappingTargetId) (entity.MappingTarget, error)
	GetTargetsByContainerId(
		containerId valueObject.ContainerId,
	) ([]entity.MappingTarget, error)
	FindOne(
		hostname *valueObject.Fqdn,
		publicPort valueObject.NetworkPort,
		protocol valueObject.NetworkProtocol,
	) (entity.Mapping, error)
}
