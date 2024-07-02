package infra

import (
	"testing"

	testHelpers "github.com/speedianet/control/src/devUtils"
	"github.com/speedianet/control/src/domain/valueObject"
)

func TestLicenseCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	transientDbSvc := testHelpers.GetTransientDbSvc()
	licenseCmdRepo := NewLicenseCmdRepo(persistentDbSvc, transientDbSvc)

	t.Run("GenerateIntegrityHash", func(t *testing.T) {
		licenseQueryRepo := NewLicenseQueryRepo(persistentDbSvc, transientDbSvc)
		licenseInfo, err := licenseQueryRepo.Read()
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}

		_, err = licenseCmdRepo.GenerateIntegrityHash(licenseInfo)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("GenerateNonceHash", func(t *testing.T) {
		_, err := licenseCmdRepo.GenerateNonceHash()
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("GenerateFingerprint", func(t *testing.T) {
		_, err := licenseCmdRepo.generateFingerprint()
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdateIntegrityHash", func(t *testing.T) {
		err := licenseCmdRepo.updateIntegrityHash()
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("Refresh", func(t *testing.T) {
		err := licenseCmdRepo.Refresh()
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdateStatus", func(t *testing.T) {
		licenseStatus, _ := valueObject.NewLicenseStatus("ACTIVE")

		err := licenseCmdRepo.UpdateStatus(licenseStatus)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("IncrementErrorCount", func(t *testing.T) {
		err := licenseCmdRepo.IncrementErrorCount()
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("ResetErrorCount", func(t *testing.T) {
		err := licenseCmdRepo.ResetErrorCount()
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})
}
