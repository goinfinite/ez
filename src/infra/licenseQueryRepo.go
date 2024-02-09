package infra

import (
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/infra/db"
)

type LicenseQueryRepo struct {
	dbSvc *db.DatabaseService
}

func NewLicenseQueryRepo(dbSvc *db.DatabaseService) LicenseQueryRepo {
	return LicenseQueryRepo{dbSvc: dbSvc}
}

func (repo LicenseQueryRepo) Get() (entity.LicenseInfo, error) {
	return entity.LicenseInfo{}, nil
}

func (repo LicenseQueryRepo) GetErrorCount() (int, error) {
	return 0, nil
}
