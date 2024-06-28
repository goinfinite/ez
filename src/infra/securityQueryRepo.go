package infra

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/infra/db"
)

type SecurityQueryRepo struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewSecurityQueryRepo(
	persistentDbSvc *db.PersistentDatabaseService,
) *SecurityQueryRepo {
	return &SecurityQueryRepo{persistentDbSvc: persistentDbSvc}
}

func (repo *SecurityQueryRepo) GetEvents(
	getDto dto.GetSecurityEvents,
) ([]entity.SecurityEvent, error) {
	return nil, nil
}
