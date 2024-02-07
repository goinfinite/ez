package infra

import (
	"time"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
)

type LicenseQueryRepo struct {
	dbSvc *db.DatabaseService
}

func NewLicenseQueryRepo(dbSvc *db.DatabaseService) *LicenseQueryRepo {
	return &LicenseQueryRepo{dbSvc: dbSvc}
}

func (repo LicenseQueryRepo) GetStatus() (entity.LicenseStatus, error) {
	licenseMethod, _ := valueObject.NewLicenseMethod("ip")
	licenseStatus, _ := valueObject.NewLicenseStatus("active")
	licenseFingerprint, _ := valueObject.NewLicenseFingerprint("fingerprint")

	todayDate := time.Now()
	expiresAt := valueObject.UnixTime(todayDate.AddDate(1, 0, 0).Unix())
	lastCheckAt := valueObject.UnixTime(todayDate.Unix())

	return entity.NewLicenseStatus(
		licenseMethod,
		licenseStatus,
		licenseFingerprint,
		expiresAt,
		lastCheckAt,
	), nil
}
