package infra

import (
	"errors"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/infra/db"
	dbModel "github.com/speedianet/control/src/infra/db/model"
)

type LicenseQueryRepo struct {
	dbSvc *db.DatabaseService
}

func NewLicenseQueryRepo(dbSvc *db.DatabaseService) LicenseQueryRepo {
	return LicenseQueryRepo{dbSvc: dbSvc}
}

func (repo LicenseQueryRepo) Get() (entity.LicenseInfo, error) {
	var licenseInfo entity.LicenseInfo

	var licenseInfoModel dbModel.LicenseInfo
	queryResult := repo.dbSvc.Orm.
		Where("id = ?", 1).
		Limit(1).
		Find(&licenseInfoModel)
	if queryResult.Error != nil {
		return licenseInfo, queryResult.Error
	}

	if queryResult.RowsAffected == 0 {
		licenseCmdRepo := NewLicenseCmdRepo(repo.dbSvc)
		err := licenseCmdRepo.Refresh()
		if err != nil {
			return licenseInfo, err
		}
	}

	queryResult = repo.dbSvc.Orm.
		Where("id = ?", 1).
		Limit(1).
		Find(&licenseInfoModel)
	if queryResult.Error != nil {
		return licenseInfo, queryResult.Error
	}

	if queryResult.RowsAffected == 0 {
		return licenseInfo, errors.New("GetLicenseInfoFailedRepeatedly")
	}

	return licenseInfoModel.ToEntity()
}

func (repo LicenseQueryRepo) GetErrorCount() (int, error) {
	return 0, nil
}
