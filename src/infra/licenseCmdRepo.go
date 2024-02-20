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
	"gorm.io/gorm"
)

type LicenseCmdRepo struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewLicenseCmdRepo(persistentDbSvc *db.PersistentDatabaseService) LicenseCmdRepo {
	return LicenseCmdRepo{persistentDbSvc: persistentDbSvc}
}

func (repo LicenseCmdRepo) Refresh() error {
	speediaApiUrl := "https://app.speedia.net/api/v1"
	apiEndpoint := "/store/product/license/verify/1/"

	licenseQueryRepo := NewLicenseQueryRepo(repo.persistentDbSvc)
	freshLicenseFingerprint, err := licenseQueryRepo.GetLicenseFingerprint()
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

	httpResponse, err := http.Get(speediaApiUrl + apiEndpoint)
	if err != nil {
		return errors.New("GetLicenseInfoFailed")
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != 200 {
		return errors.New("LicenseInfoBadStatusCode")
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

	rawExpiresAt, assertOk := parsedBody["expiresAt"].(float64)
	if !assertOk {
		return errors.New("ParseLicenseExpiresAtFailed")
	}
	expiresAt := valueObject.UnixTime(int64(rawExpiresAt))

	todayDate := time.Now()
	lastCheckAt := valueObject.UnixTime(todayDate.Unix())
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
	return repo.persistentDbSvc.Handler.Save(&licenseInfoModel).Error
}

func (repo LicenseCmdRepo) UpdateStatus(status valueObject.LicenseStatus) error {
	licenseInfoModel := dbModel.LicenseInfo{}
	updateResult := repo.persistentDbSvc.Handler.Model(&licenseInfoModel).
		Where("id = ?", 1).
		Update("status", status.String())
	return updateResult.Error
}

func (repo LicenseCmdRepo) IncrementErrorCount() error {
	licenseInfoModel := dbModel.LicenseInfo{}
	updateResult := repo.persistentDbSvc.Handler.Model(&licenseInfoModel).
		Where("id = ?", 1).
		UpdateColumn("error_count", gorm.Expr("error_count + ?", 1))
	return updateResult.Error
}

func (repo LicenseCmdRepo) ResetErrorCount() error {
	licenseInfoModel := dbModel.LicenseInfo{}
	updateResult := repo.persistentDbSvc.Handler.Model(&licenseInfoModel).
		Where("id = ?", 1).
		Update("error_count", 0)
	return updateResult.Error
}
