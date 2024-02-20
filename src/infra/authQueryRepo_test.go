package infra

import (
	"os"
	"testing"
	"time"

	testHelpers "github.com/speedianet/control/src/devUtils"
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/valueObject"
)

func TestAuthQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistDbSvc := testHelpers.GetPersistentDbSvc()
	authQueryRepo := NewAuthQueryRepo(persistDbSvc)
	accCmdRepo := NewAccCmdRepo(persistDbSvc)

	t.Run("ValidLoginCredentials", func(t *testing.T) {
		login := dto.NewLogin(
			valueObject.NewUsernamePanic(os.Getenv("DUMMY_USER_NAME")),
			valueObject.NewPasswordPanic(os.Getenv("DUMMY_USER_PASS")),
		)
		isValid := authQueryRepo.IsLoginValid(login)
		if !isValid {
			t.Error("Expected valid login credentials, but got invalid")
		}
	})

	t.Run("InvalidLoginCredentials", func(t *testing.T) {
		login := dto.NewLogin(
			valueObject.NewUsernamePanic(os.Getenv("DUMMY_USER_NAME")),
			valueObject.NewPasswordPanic("wrongPassword"),
		)
		isValid := authQueryRepo.IsLoginValid(login)
		if isValid {
			t.Error("Expected invalid login credentials, but got valid")
		}
	})

	t.Run("ValidSessionAccessToken", func(t *testing.T) {
		authCmdRepo := AuthCmdRepo{}

		token, _ := authCmdRepo.GenerateSessionToken(
			valueObject.AccountId(1000),
			valueObject.UnixTime(
				time.Now().Add(3*time.Hour).Unix(),
			),
			valueObject.NewIpAddressPanic("127.0.0.1"),
		)

		_, err := authQueryRepo.GetAccessTokenDetails(token.TokenStr)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("InvalidSessionAccessToken", func(t *testing.T) {
		invalidToken := valueObject.NewAccessTokenStrPanic(
			"invalidTokenInvalidTokenInvalidTokenInvalidTokenInvalidToken",
		)
		_, err := authQueryRepo.GetAccessTokenDetails(invalidToken)
		if err == nil {
			t.Error("ExpectingError")
		}
	})

	t.Run("ValidAccountApiKey", func(t *testing.T) {
		apiKey, err := accCmdRepo.UpdateApiKey(
			valueObject.NewAccountIdPanic(os.Getenv("DUMMY_USER_ID")),
		)
		if err != nil {
			t.Error(err)
		}

		_, err = authQueryRepo.GetAccessTokenDetails(apiKey)
		if err != nil {
			t.Error(err)
		}
	})
}
