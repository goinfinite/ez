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

	t.Run("Get", func(t *testing.T) {
		_, err := licenseQueryRepo.Get()
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("GetIntegrityHash", func(t *testing.T) {
		_, err := licenseQueryRepo.GetIntegrityHash()
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})
}
