package infra

import (
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
)

type MappingQueryRepo struct {
	dbSvc *db.DatabaseService
}

func NewMappingQueryRepo(dbSvc *db.DatabaseService) *MappingQueryRepo {
	return &MappingQueryRepo{dbSvc: dbSvc}
}

func (repo MappingQueryRepo) Get() ([]entity.Mapping, error) {
	return []entity.Mapping{}, nil
}

func (repo MappingQueryRepo) GetById(id valueObject.MappingId) (entity.Mapping, error) {
	return entity.Mapping{}, nil
}

func (repo MappingQueryRepo) GetByHostPortProtocol(
	hostname *valueObject.Fqdn,
	port valueObject.NetworkPort,
	protocol valueObject.NetworkProtocol,
) (entity.Mapping, error) {
	return entity.Mapping{}, nil
}
