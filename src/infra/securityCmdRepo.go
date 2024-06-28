package infra

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/infra/db"
)

type SecurityCmdRepo struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewSecurityCmdRepo(
	persistentDbSvc *db.PersistentDatabaseService,
) *SecurityCmdRepo {
	return &SecurityCmdRepo{persistentDbSvc: persistentDbSvc}
}

func (repo *SecurityCmdRepo) CreateEvent(createDto dto.CreateSecurityEvent) error {
	return nil
}

func (repo *SecurityCmdRepo) DeleteEvents(deleteDto dto.DeleteSecurityEvents) error {
	return nil
}
