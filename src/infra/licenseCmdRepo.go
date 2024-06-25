package infra

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
	dbModel "github.com/speedianet/control/src/infra/db/model"
	infraHelper "github.com/speedianet/control/src/infra/helper"
	"gorm.io/gorm"
)

type LicenseCmdRepo struct {
	persistentDbSvc *db.PersistentDatabaseService
	transientDbSvc  *db.TransientDatabaseService
}

func NewLicenseCmdRepo(
	persistentDbSvc *db.PersistentDatabaseService,
	transientDbSvc *db.TransientDatabaseService,
) *LicenseCmdRepo {
	return &LicenseCmdRepo{
		persistentDbSvc: persistentDbSvc,
		transientDbSvc:  transientDbSvc,
	}
}

func (repo *LicenseCmdRepo) GenerateIntegrityHash(
	licenseInfo entity.LicenseInfo,
) (valueObject.Hash, error) {
	var integrityHash valueObject.Hash

	licenseInfoJson, err := json.Marshal(licenseInfo)
	if err != nil {
		return integrityHash, errors.New("MarshalLicenseInfoFailed: " + err.Error())
	}

	licenseInfoHashStr := infraHelper.GenHash(string(licenseInfoJson))
	return valueObject.NewHash(licenseInfoHashStr)
}

func (repo *LicenseCmdRepo) GenerateNonceHash() (valueObject.Hash, error) {
	var nonceHash valueObject.Hash

	currentHourInEpoch, err := infraHelper.RunCmdWithSubShell(
		"date -d \"$(date +'%Y-%m-%d %H:00:00')\" +%s",
	)
	if err != nil {
		return nonceHash, err
	}

	nonceHashStr := infraHelper.GenShortHash(currentHourInEpoch)

	return valueObject.NewHash(nonceHashStr)
}

func (repo *LicenseCmdRepo) generateFingerprint() (
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

	publicIp, err := infraHelper.GetPublicIpAddress(repo.transientDbSvc)
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

func (repo *LicenseCmdRepo) updateIntegrityHash() error {
	licenseQueryRepo := NewLicenseQueryRepo(repo.persistentDbSvc, repo.transientDbSvc)

	licenseInfo, err := licenseQueryRepo.Get()
	if err != nil {
		return errors.New("GetLicenseInfoFailed: " + err.Error())
	}

	licenseInfoHash, err := repo.GenerateIntegrityHash(licenseInfo)
	if err != nil {
		return errors.New("GenerateLicenseInfoHashFailed: " + err.Error())
	}

	err = repo.transientDbSvc.Set(LicenseInfoHashKey, licenseInfoHash.String())
	if err != nil {
		return errors.New("SetLicenseInfoHashFailed: " + err.Error())
	}

	return nil
}

func (repo *LicenseCmdRepo) Refresh() error {
	speediaApiUrl := "https://app.speedia.net/api/v1"
	apiEndpoint := "/store/product/license/verify/1/"

	freshLicenseFingerprint, err := repo.generateFingerprint()
	if err != nil {
		return errors.New("GetLicenseFingerprintFailed")
	}

	apiEndpoint += "?fingerprint=" + freshLicenseFingerprint.String()

	licenseMethod, _ := valueObject.NewLicenseMethod("ip")

	keyStr := os.Getenv("LICENSE_KEY")
	if keyStr != "" {
		licenseMethod, _ = valueObject.NewLicenseMethod("key")
		apiEndpoint += "&key=" + keyStr
	}

	httpRequest, err := http.NewRequest(http.MethodGet, speediaApiUrl+apiEndpoint, nil)
	if err != nil {
		return errors.New("LicenseServerRequestError: " + err.Error())
	}
	httpRequest.Header.Set("User-Agent", "Speedia Control/1.0")

	httpClient := &http.Client{
		Timeout: time.Second *
			15,
	}

	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		return errors.New("LicenseServerResponseError: " + err.Error())
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != 200 {
		return errors.New("LicenseServerBadStatusCode")
	}

	responseBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return errors.New("ReadLicenseInfoFailed")
	}

	var parsedResponse map[string]interface{}
	err = json.Unmarshal(responseBody, &parsedResponse)
	if err != nil {
		return errors.New("ParseLicenseInfoFailed")
	}

	parsedBody, assertOk := parsedResponse["body"].(map[string]interface{})
	if !assertOk {
		return errors.New("ParseLicenseBodyFailed")
	}

	rawStatus, assertOk := parsedBody["status"].(string)
	if !assertOk {
		return errors.New("ParseLicenseStatusFailed")
	}
	licenseStatus, err := valueObject.NewLicenseStatus(rawStatus)
	if err != nil {
		return err
	}

	rawLicenseFingerprint, assertOk := parsedBody["licenseFingerprint"].(string)
	if !assertOk {
		return errors.New("ParseLicenseFingerprintFailed")
	}
	licenseFingerprint, err := valueObject.NewLicenseFingerprint(rawLicenseFingerprint)
	if err != nil {
		return err
	}

	expiresAt, err := valueObject.NewUnixTime(parsedBody["expiresAt"])
	if err != nil {
		return errors.New("ParseLicenseExpiresAtFailed")
	}

	lastCheckAt := valueObject.NewUnixTimeNow()
	updatedAt := lastCheckAt

	licenseInfoEntity := entity.NewLicenseInfo(
		licenseMethod,
		licenseStatus,
		licenseFingerprint,
		0,
		expiresAt,
		lastCheckAt,
		updatedAt,
	)

	licenseInfoModel := dbModel.LicenseInfo{}.ToModel(licenseInfoEntity)
	err = repo.persistentDbSvc.Handler.Save(&licenseInfoModel).Error
	if err != nil {
		return errors.New("SaveLicenseInfoFailed: " + err.Error())
	}

	return repo.updateIntegrityHash()
}

func (repo *LicenseCmdRepo) UpdateStatus(status valueObject.LicenseStatus) error {
	licenseInfoModel := dbModel.LicenseInfo{}

	err := repo.persistentDbSvc.Handler.Model(&licenseInfoModel).
		Where("id = ?", 1).
		Update("status", status.String()).Error

	if err != nil {
		return err
	}

	return repo.updateIntegrityHash()
}

func (repo *LicenseCmdRepo) IncrementErrorCount() error {
	licenseInfoModel := dbModel.LicenseInfo{}

	err := repo.persistentDbSvc.Handler.Model(&licenseInfoModel).
		Where("id = ?", 1).
		UpdateColumn("error_count", gorm.Expr("error_count + ?", 1)).Error
	if err != nil {
		return err
	}

	return repo.updateIntegrityHash()
}

func (repo *LicenseCmdRepo) ResetErrorCount() error {
	licenseInfoModel := dbModel.LicenseInfo{}

	err := repo.persistentDbSvc.Handler.Model(&licenseInfoModel).
		Where("id = ?", 1).
		Update("error_count", 0).Error
	if err != nil {
		return err
	}

	return repo.updateIntegrityHash()
}
