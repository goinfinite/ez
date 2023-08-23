package infra

import (
	"os"
	"testing"

	testHelpers "github.com/speedianet/sfm/src/devUtils"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

func TestAccQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	authQueryRepo := AccQueryRepo{}

	t.Run("GetValidAccounts", func(t *testing.T) {
		_, err := authQueryRepo.Get()
		if err != nil {
			t.Error("UnexpectedError")
		}
	})

	t.Run("GetValidAccountByUsername", func(t *testing.T) {
		username := valueObject.NewUsernamePanic(os.Getenv("DUMMY_USER_NAME"))

		_, err := authQueryRepo.GetByUsername(username)
		if err != nil {
			t.Error("UnexpectedError")
		}
	})

	t.Run("GetValidAccountById", func(t *testing.T) {
		accountId := valueObject.NewAccountIdFromStringPanic(os.Getenv("DUMMY_USER_ID"))

		_, err := authQueryRepo.GetById(accountId)
		if err != nil {
			t.Error("UnexpectedError")
		}
	})

	t.Run("GetInvalidAccount", func(t *testing.T) {
		username := valueObject.NewUsernamePanic("invalid")

		_, err := authQueryRepo.GetByUsername(username)
		if err == nil {
			t.Error("ExpectingError")
		}
	})
}
