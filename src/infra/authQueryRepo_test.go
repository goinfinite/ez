package infra

import (
	"os"
	"testing"
	"time"

	testHelpers "github.com/speedianet/sfm/src/devUtils"
	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

func TestAuthQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()

	t.Run("ValidLoginCredentials", func(t *testing.T) {
		login := dto.NewLogin(
			valueObject.NewUsernamePanic(os.Getenv("DUMMY_USER_NAME")),
			valueObject.NewPasswordPanic(os.Getenv("DUMMY_USER_PASS")),
		)
		authQueryRepo := AuthQueryRepo{}
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
		authQueryRepo := AuthQueryRepo{}
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

		authQueryRepo := AuthQueryRepo{}
		_, err := authQueryRepo.GetAccessTokenDetails(token.TokenStr)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("InvalidSessionAccessToken", func(t *testing.T) {
		authQueryRepo := AuthQueryRepo{}
		invalidToken := valueObject.NewAccessTokenStrPanic(
			"invalidTokenInvalidTokenInvalidTokenInvalidTokenInvalidToken",
		)
		_, err := authQueryRepo.GetAccessTokenDetails(invalidToken)
		if err == nil {
			t.Error("ExpectingError")
		}
	})

	t.Run("ValidAccountApiKey", func(t *testing.T) {
		accCmdRepo := AccCmdRepo{}

		apiKey, err := accCmdRepo.UpdateApiKey(
			valueObject.NewAccountIdPanic(os.Getenv("DUMMY_USER_ID")),
		)
		if err != nil {
			t.Error(err)
		}

		authQueryRepo := AuthQueryRepo{}
		_, err = authQueryRepo.GetAccessTokenDetails(apiKey)
		if err != nil {
			t.Error(err)
		}
	})
}
