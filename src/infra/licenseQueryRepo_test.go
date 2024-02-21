package infra

import (
	"testing"

	testHelpers "github.com/speedianet/control/src/devUtils"
)

func TestLicenseQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	transientDbSbc := testHelpers.GetTransientDbSvc()
	licenseQueryRepo := NewLicenseQueryRepo(persistentDbSvc, transientDbSbc)

	t.Run("GetLicenseInfo", func(t *testing.T) {
		_, err := licenseQueryRepo.Get()
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("GetNonceHash", func(t *testing.T) {
		_, err := licenseQueryRepo.GetNonceHash()
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("GetLicenseFingerprint", func(t *testing.T) {
		_, err := licenseQueryRepo.GetFingerprint()
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})
}
