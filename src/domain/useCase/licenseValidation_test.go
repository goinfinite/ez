package useCase

import (
	"testing"

	testHelpers "github.com/goinfinite/ez/src/devUtils"
	"github.com/goinfinite/ez/src/infra"
	infraHelper "github.com/goinfinite/ez/src/infra/helper"
)

func TestLicenseValidation(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	transientDbSvc := testHelpers.GetTransientDbSvc()
	licenseQueryRepo := infra.NewLicenseQueryRepo(persistentDbSvc, transientDbSvc)
	licenseCmdRepo := infra.NewLicenseCmdRepo(persistentDbSvc, transientDbSvc)

	t.Run("LicenseValidationWithPerfectConditions", func(t *testing.T) {
		err := LicenseValidation(licenseQueryRepo, licenseCmdRepo)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("LicenseValidationWithLicenseServerUnreachable", func(t *testing.T) {
		_, err := infraHelper.RunCmdWithSubShell(
			"echo \"127.0.0.1 app.speedia.net\" >> /etc/hosts",
		)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}

		err = LicenseValidation(licenseQueryRepo, licenseCmdRepo)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}

		licenseInfo, err := licenseQueryRepo.Read()
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}

		if licenseInfo.ErrorCount < 1 {
			t.Errorf("UnexpectedErrorCount: %v", licenseInfo.ErrorCount)
		}

		_, err = infraHelper.RunCmdWithSubShell(
			"sed -i '$ d' /etc/hosts",
		)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})
}
