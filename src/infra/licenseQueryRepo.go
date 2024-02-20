package infra

import (
	"errors"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
	dbModel "github.com/speedianet/control/src/infra/db/model"
	infraHelper "github.com/speedianet/control/src/infra/helper"
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

func (repo LicenseQueryRepo) GetErrorCount() (uint, error) {
	return 0, nil
}

func (repo LicenseQueryRepo) GetLicenseFingerprint() (
	valueObject.LicenseFingerprint,
	error,
) {
	var fingerprint valueObject.LicenseFingerprint

	hwUuid, err := infraHelper.RunCmdWithSubShell("dmidecode -t system | awk '/UUID/{print $2}'")
	if err != nil {
		return fingerprint, err
	}

	rootFsUuid, err := infraHelper.RunCmdWithSubShell(
		"blkid $(df --output=source / | tail -1) | grep -oP '(?<= UUID=\")[a-fA-F0-9-]+'",
	)
	if err != nil {
		return fingerprint, err
	}

	privateIp, err := infraHelper.RunCmdWithSubShell("hostname -I | awk '{print $1}'")
	if err != nil {
		return fingerprint, err
	}

	installationUnixTime, err := infraHelper.RunCmdWithSubShell(
		"stat -c %W " + db.DatabaseFilePath,
	)
	if err != nil {
		return fingerprint, err
	}

	fingerprintFirstPart := hwUuid + rootFsUuid + privateIp + installationUnixTime
	firstPartShortHashStr := infraHelper.GenShortHash(fingerprintFirstPart)

	publicIp, err := infraHelper.GetPublicIpAddress()
	if err != nil {
		return fingerprint, err
	}

	macAddress, err := infraHelper.RunCmdWithSubShell(
		"ip link show | awk '/link\\/ether/{print $2}'",
	)
	if err != nil {
		return fingerprint, err
	}

	fingerprintSecondPart := publicIp.String() + macAddress
	secondPartShortHashStr := infraHelper.GenShortHash(fingerprintSecondPart)

	currentHourInEpoch, err := infraHelper.RunCmdWithSubShell(
		"date -d \"$(date +'%Y-%m-%d %H:00:00')\" +%s",
	)
	if err != nil {
		return fingerprint, err
	}

	fingerprintThirdPart := currentHourInEpoch
	thirdPartShortHashStr := infraHelper.GenShortHash(fingerprintThirdPart)

	return valueObject.NewLicenseFingerprint(
		firstPartShortHashStr + "-" + secondPartShortHashStr + "-" + thirdPartShortHashStr,
	)
}
