package infra

import (
	"testing"

	testHelpers "github.com/speedianet/control/src/devUtils"
)

func TestLicenseQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	dbSvc := testHelpers.GetDbSvc()
	licenseQueryRepo := NewLicenseQueryRepo(dbSvc)

	t.Run("GetLicenseStatus", func(t *testing.T) {
		_, err := licenseQueryRepo.GetStatus()
		if err != nil {
			t.Error("UnexpectedError")
		}
	})
}
