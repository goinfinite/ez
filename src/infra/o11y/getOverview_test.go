package o11yInfra

import (
	"testing"

	testHelpers "github.com/speedianet/control/src/devUtils"
)

func TestGetOverview(t *testing.T) {
	testHelpers.LoadEnvVars()
	transientDbSvc := testHelpers.GetTransientDbSvc()
	o11yQueryRepo := NewO11yQueryRepo(transientDbSvc)

	t.Run("GetOverview", func(t *testing.T) {
		_, err := o11yQueryRepo.GetOverview()
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})
}
