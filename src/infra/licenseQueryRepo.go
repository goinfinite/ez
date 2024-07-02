package infra

import (
	"errors"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
	dbModel "github.com/speedianet/control/src/infra/db/model"
)

const (
	LicenseInfoHashKey = "licenseInfoHash"
)

type LicenseQueryRepo struct {
	persistentDbSvc *db.PersistentDatabaseService
	transientDbSvc  *db.TransientDatabaseService
}

func NewLicenseQueryRepo(
	persistentDbSvc *db.PersistentDatabaseService,
	transientDbSvc *db.TransientDatabaseService,
) *LicenseQueryRepo {
	return &LicenseQueryRepo{
		persistentDbSvc: persistentDbSvc,
		transientDbSvc:  transientDbSvc,
	}
}

func (repo *LicenseQueryRepo) Read() (entity.LicenseInfo, error) {
	var licenseInfo entity.LicenseInfo

	var licenseInfoModel dbModel.LicenseInfo
	queryResult := repo.persistentDbSvc.Handler.
		Where("id = ?", 1).
		Limit(1).
		Find(&licenseInfoModel)
	if queryResult.Error != nil {
		return licenseInfo, queryResult.Error
	}

	if queryResult.RowsAffected == 0 {
		return licenseInfo, errors.New("LicenseInfoNotFound")
	}

	return licenseInfoModel.ToEntity()
}

func (repo *LicenseQueryRepo) ReadIntegrityHash() (valueObject.Hash, error) {
	var licenseInfoHash valueObject.Hash

	licenseInfoHashStr, err := repo.transientDbSvc.Get(LicenseInfoHashKey)
	if err != nil {
		return licenseInfoHash, err
	}

	licenseInfoHash, err = valueObject.NewHash(licenseInfoHashStr)
	if err != nil {
		return licenseInfoHash, err
	}

	return licenseInfoHash, nil
}
