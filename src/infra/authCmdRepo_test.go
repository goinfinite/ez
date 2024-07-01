package infra

import (
	"testing"
	"time"

	testHelpers "github.com/speedianet/control/src/devUtils"
	"github.com/speedianet/control/src/domain/valueObject"
)

func TestAuthCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()

	t.Run("GetSessionToken", func(t *testing.T) {
		authCmdRepo := AuthCmdRepo{}
		_, err := authCmdRepo.GenerateSessionToken(
			valueObject.AccountId(1000),
			valueObject.NewUnixTimeAfterNow(3*time.Hour),
			valueObject.NewLocalhostIpAddress(),
		)
		if err != nil {
			t.Error(err)
		}
	})
}
