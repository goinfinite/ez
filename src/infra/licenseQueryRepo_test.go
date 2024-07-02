package infra

import (
	"testing"

	testHelpers "github.com/speedianet/control/src/devUtils"
)

func TestLicenseQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	transientDbSvc := testHelpers.GetTransientDbSvc()
	licenseQueryRepo := NewLicenseQueryRepo(persistentDbSvc, transientDbSvc)
	licenseCmdRepo := NewLicenseCmdRepo(persistentDbSvc, transientDbSvc)

	t.Run("Get", func(t *testing.T) {
		_, err := licenseQueryRepo.Read()
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("ReadIntegrityHash", func(t *testing.T) {
		err := licenseCmdRepo.updateIntegrityHash()
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}

		_, err = licenseQueryRepo.ReadIntegrityHash()
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})
}
