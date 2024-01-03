package repository

import (
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type MappingQueryRepo interface {
	Get() ([]entity.Mapping, error)
	GetById(id valueObject.MappingId) (entity.Mapping, error)
	GetTargetById(id valueObject.MappingTargetId) (entity.MappingTarget, error)
	FindOne(
		hostname *valueObject.Fqdn,
		port valueObject.NetworkPort,
		protocol valueObject.NetworkProtocol,
	) (entity.Mapping, error)
}
