package infra

import (
	"os"
	"testing"

	testHelpers "github.com/goinfinite/fleet/src/devUtils"
	"github.com/goinfinite/fleet/src/domain/valueObject"
)

func TestAccQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	dbSvc := testHelpers.GetDbSvc()
	accQueryRepo := NewAccQueryRepo(dbSvc)

	t.Run("GetValidAccounts", func(t *testing.T) {
		_, err := accQueryRepo.Get()
		if err != nil {
			t.Error("UnexpectedError")
		}
	})

	t.Run("GetValidAccountByUsername", func(t *testing.T) {
		username := valueObject.NewUsernamePanic(os.Getenv("DUMMY_USER_NAME"))

		_, err := accQueryRepo.GetByUsername(username)
		if err != nil {
			t.Error("UnexpectedError")
		}
	})

	t.Run("GetValidAccountById", func(t *testing.T) {
		accountId := valueObject.NewAccountIdPanic(os.Getenv("DUMMY_USER_ID"))

		_, err := accQueryRepo.GetById(accountId)
		if err != nil {
			t.Error("UnexpectedError")
		}
	})

	t.Run("GetInvalidAccount", func(t *testing.T) {
		username := valueObject.NewUsernamePanic("invalid")

		_, err := accQueryRepo.GetByUsername(username)
		if err == nil {
			t.Error("ExpectingError")
		}
	})
}
