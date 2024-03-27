package infraHelper

import (
	"testing"

	testHelpers "github.com/speedianet/control/src/devUtils"
)

func TestGetPublicIpAddress(t *testing.T) {
	testHelpers.LoadEnvVars()
	transientDbSvc := testHelpers.GetTransientDbSvc()

	t.Run("GetValidPublicIpAddress", func(t *testing.T) {
		_, err := GetPublicIpAddress(transientDbSvc)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})
}
