package infra

import (
	"testing"

	testHelpers "github.com/speedianet/control/src/devUtils"
)

func TestContainerProfileQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	dbSvc := testHelpers.GetDbSvc()
	profileQueryRepo := NewContainerProfileQueryRepo(dbSvc)

	t.Run("GetContainerProfiles", func(t *testing.T) {
		_, err := profileQueryRepo.Get()
		if err != nil {
			t.Errorf("GetContainerProfilesFailed: %v", err)
		}
	})

	t.Run("GetContainerProfileById", func(t *testing.T) {
		_, err := profileQueryRepo.GetById(1)
		if err != nil {
			t.Errorf("GetContainerProfileByIdFailed: %v", err)
		}
	})

	t.Run("GetProfileByIdWithInvalidId", func(t *testing.T) {
		_, err := profileQueryRepo.GetById(100)
		if err == nil {
			t.Errorf("GetProfileByIdWithInvalidIdFailed: %v", err)
		}
	})
}
