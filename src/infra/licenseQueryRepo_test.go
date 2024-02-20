package infra

import (
	"testing"

	testHelpers "github.com/speedianet/control/src/devUtils"
)

func TestLicenseQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	dbSvc := testHelpers.GetDbSvc()
	licenseQueryRepo := NewLicenseQueryRepo(dbSvc)

	t.Run("GetLicenseInfo", func(t *testing.T) {
		_, err := licenseQueryRepo.Get()
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("GetErrorCount", func(t *testing.T) {
		_, err := licenseQueryRepo.GetErrorCount()
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("GetLicenseFingerprint", func(t *testing.T) {
		_, err := licenseQueryRepo.GetLicenseFingerprint()
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})
}
