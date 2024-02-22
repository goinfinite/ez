package useCase

import (
	"testing"

	testHelpers "github.com/speedianet/control/src/devUtils"
	"github.com/speedianet/control/src/infra"
)

func TestLicenseValidation(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	transientDbSvc := testHelpers.GetTransientDbSvc()
	licenseQueryRepo := infra.NewLicenseQueryRepo(persistentDbSvc, transientDbSvc)
	licenseCmdRepo := infra.NewLicenseCmdRepo(persistentDbSvc, transientDbSvc)

	t.Run("LicenseValidation", func(t *testing.T) {
		err := LicenseValidation(licenseQueryRepo, licenseCmdRepo)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})
}
