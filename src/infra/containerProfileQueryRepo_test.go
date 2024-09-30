package infra

import (
	"testing"

	testHelpers "github.com/goinfinite/ez/src/devUtils"
)

func TestContainerProfileQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	profileQueryRepo := NewContainerProfileQueryRepo(persistentDbSvc)

	t.Run("ReadContainerProfiles", func(t *testing.T) {
		_, err := profileQueryRepo.Read()
		if err != nil {
			t.Errorf("ReadContainerProfilesFailed: %v", err)
		}
	})

	t.Run("GetContainerProfileById", func(t *testing.T) {
		_, err := profileQueryRepo.ReadById(1)
		if err != nil {
			t.Errorf("GetContainerProfileByIdFailed: %v", err)
		}
	})

	t.Run("GetProfileByIdWithInvalidId", func(t *testing.T) {
		_, err := profileQueryRepo.ReadById(100)
		if err == nil {
			t.Errorf("GetProfileByIdWithInvalidIdFailed: %v", err)
		}
	})
}
