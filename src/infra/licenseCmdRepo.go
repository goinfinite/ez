package infra

import (
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
)

type LicenseCmdRepo struct {
	dbSvc *db.DatabaseService
}

func NewLicenseCmdRepo(dbSvc *db.DatabaseService) LicenseCmdRepo {
	return LicenseCmdRepo{dbSvc: dbSvc}
}

func (repo LicenseCmdRepo) RefreshStatus() error {
	return nil
}

func (repo LicenseCmdRepo) UpdateStatus(valueObject.LicenseStatus) error {
	return nil
}

func (repo LicenseCmdRepo) IncrementErrorCount() error {
	return nil
}

func (repo LicenseCmdRepo) ResetErrorCount() error {
	return nil
}
