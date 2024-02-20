package infra

import (
	"testing"

	testHelpers "github.com/speedianet/control/src/devUtils"
	"github.com/speedianet/control/src/domain/valueObject"
)

func TestLicenseCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistDbSvc := testHelpers.GetPersistentDbSvc()
	licenseCmdRepo := NewLicenseCmdRepo(persistDbSvc)

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
