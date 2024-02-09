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
)

type LicenseQueryRepo struct {
	dbSvc *db.DatabaseService
}

func NewLicenseQueryRepo(dbSvc *db.DatabaseService) LicenseQueryRepo {
	return LicenseQueryRepo{dbSvc: dbSvc}
}

func (repo LicenseQueryRepo) Get() (entity.LicenseInfo, error) {
	var licenseInfo entity.LicenseInfo

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
		return licenseInfo, errors.New("GetLicenseInfoFailed")
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != 200 {
		return licenseInfo, errors.New("LicenseInfoBadStatusCode")
	}

	responseBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return licenseInfo, errors.New("ReadLicenseInfoFailed")
	}

	var parsedResponse map[string]interface{}
	err = json.Unmarshal(responseBody, &parsedResponse)
	if err != nil {
		return licenseInfo, errors.New("ParseLicenseInfoFailed")
	}

	parsedBody, assertOk := parsedResponse["body"].(map[string]interface{})
	if !assertOk {
		return licenseInfo, errors.New("ParseLicenseBodyFailed")
	}

	rawStatus, assertOk := parsedBody["status"].(string)
	if !assertOk {
		return licenseInfo, errors.New("ParseLicenseStatusFailed")
	}
	licenseStatus, err := valueObject.NewLicenseStatus(rawStatus)
	if err != nil {
		return licenseInfo, err
	}

	// TODO: Implement license fingerprint
	licenseFingerprint, _ := valueObject.NewLicenseFingerprint("fingerprint")

	rawExpiresAt, assertOk := parsedBody["expiresAt"].(float64)
	if !assertOk {
		return licenseInfo, errors.New("ParseLicenseExpiresAtFailed")
	}
	expiresAt := valueObject.UnixTime(int64(rawExpiresAt))

	todayDate := time.Now()
	lastCheckAt := valueObject.UnixTime(todayDate.Unix())

	return entity.NewLicenseInfo(
		licenseMethod,
		licenseStatus,
		licenseFingerprint,
		expiresAt,
		lastCheckAt,
	), nil
}

func (repo LicenseQueryRepo) GetErrorCount() (int, error) {
	return 0, nil
}
