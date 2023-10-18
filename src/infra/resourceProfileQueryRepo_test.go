package infra

import (
	"testing"

	testHelpers "github.com/speedianet/sfm/src/devUtils"
)

func TestResourceProfileQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	dbSvc := testHelpers.GetDbSvc()
	profileQueryRepo := NewResourceProfileQueryRepo(dbSvc)

	t.Run("GetResourceProfiles", func(t *testing.T) {
		_, err := profileQueryRepo.Get()
		if err != nil {
			t.Errorf("GetResourceProfilesFailed: %v", err)
		}
	})
}
