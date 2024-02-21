package infra

import (
	"errors"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
	dbModel "github.com/speedianet/control/src/infra/db/model"
	infraHelper "github.com/speedianet/control/src/infra/helper"
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

func (repo LicenseQueryRepo) Get() (entity.LicenseInfo, error) {
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

func (repo LicenseQueryRepo) GetIntegrityHash() (valueObject.Hash, error) {
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

func (repo LicenseQueryRepo) GetFingerprint() (
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

	licenseCmdRepo := NewLicenseCmdRepo(repo.persistentDbSvc, repo.transientDbSvc)
	thirdPartShortHashStr, err := licenseCmdRepo.GenerateNonceHash()
	if err != nil {
		return fingerprint, err
	}

	return valueObject.NewLicenseFingerprint(
		firstPartShortHashStr + "-" +
			secondPartShortHashStr + "-" +
			thirdPartShortHashStr.String(),
	)
}
