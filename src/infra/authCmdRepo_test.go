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
			valueObject.UnixTime(
				time.Now().Add(3*time.Hour).Unix(),
			),
			valueObject.NewIpAddressPanic("127.0.0.1"),
		)
		if err != nil {
			t.Error(err)
		}
	})
}
