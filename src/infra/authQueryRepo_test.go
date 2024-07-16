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
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	authQueryRepo := NewAuthQueryRepo(persistentDbSvc)
	accountCmdRepo := NewAccountCmdRepo(persistentDbSvc)

	t.Run("ValidLoginCredentials", func(t *testing.T) {
		username, _ := valueObject.NewUsername(os.Getenv("DUMMY_USER_NAME"))
		password, _ := valueObject.NewPassword(os.Getenv("DUMMY_USER_PASS"))

		login := dto.NewLogin(username, password, nil)
		isValid := authQueryRepo.IsLoginValid(login)
		if !isValid {
			t.Error("Expected valid login credentials, but got invalid")
		}
	})

	t.Run("InvalidLoginCredentials", func(t *testing.T) {
		username, _ := valueObject.NewUsername(os.Getenv("DUMMY_USER_NAME"))
		password, _ := valueObject.NewPassword("wrongPassword")

		login := dto.NewLogin(username, password, nil)
		isValid := authQueryRepo.IsLoginValid(login)
		if isValid {
			t.Error("Expected invalid login credentials, but got valid")
		}
	})

	t.Run("ValidSessionAccessToken", func(t *testing.T) {
		authCmdRepo := AuthCmdRepo{}

		token, _ := authCmdRepo.CreateSessionToken(
			valueObject.AccountId(1000),
			valueObject.NewUnixTimeAfterNow(3*time.Hour),
			valueObject.NewLocalhostIpAddress(),
		)

		_, err := authQueryRepo.ReadAccessTokenDetails(token.TokenStr)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("InvalidSessionAccessToken", func(t *testing.T) {
		invalidToken, _ := valueObject.NewAccessTokenValue(
			"invalidTokenInvalidTokenInvalidTokenInvalidTokenInvalidToken",
		)
		_, err := authQueryRepo.ReadAccessTokenDetails(invalidToken)
		if err == nil {
			t.Error("ExpectingError")
		}
	})

	t.Run("ValidAccountApiKey", func(t *testing.T) {
		apiKey, err := accountCmdRepo.UpdateApiKey(
			valueObject.NewAccountIdPanic(os.Getenv("DUMMY_USER_ID")),
		)
		if err != nil {
			t.Error(err)
		}

		_, err = authQueryRepo.ReadAccessTokenDetails(apiKey)
		if err != nil {
			t.Error(err)
		}
	})
}
