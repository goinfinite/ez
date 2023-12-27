package infra

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
)

type MappingCmdRepo struct {
	dbSvc *db.DatabaseService
}

func NewMappingCmdRepo(dbSvc *db.DatabaseService) *MappingCmdRepo {
	return &MappingCmdRepo{dbSvc: dbSvc}
}

func (repo MappingCmdRepo) Add(mapping dto.AddMapping) error {
	return nil
}

func (repo MappingCmdRepo) Delete(id valueObject.MappingId) error {
	return nil
}
