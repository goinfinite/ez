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
)

type LicenseCmdRepo struct {
	dbSvc *db.DatabaseService
}

func NewLicenseCmdRepo(dbSvc *db.DatabaseService) LicenseCmdRepo {
	return LicenseCmdRepo{dbSvc: dbSvc}
}

func (repo LicenseCmdRepo) Refresh() error {
	speediaApiUrl := "https://app.speedia.net/api/v1"
	apiEndpoint := "/store/product/license/verify/1/"

	licenseMethod, _ := valueObject.NewLicenseMethod("ip")

	keyStr := os.Getenv("LICENSE_KEY")
	if keyStr != "" {
		licenseMethod, _ = valueObject.NewLicenseMethod("key")
		apiEndpoint += "?key=" + keyStr
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

	licenseQueryRepo := NewLicenseQueryRepo(repo.dbSvc)
	licenseFingerprint, err := licenseQueryRepo.GetLicenseFingerprint()
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
		expiresAt,
		lastCheckAt,
		updatedAt,
	)

	licenseInfoModel := dbModel.LicenseInfo{}.ToModel(licenseInfoEntity)
	return repo.dbSvc.Orm.Save(&licenseInfoModel).Error
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
