package infra

import (
	"testing"

	testHelpers "github.com/speedianet/control/src/devUtils"
)

func TestContainerImageQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	containerImageQueryRepo := NewContainerImageQueryRepo(persistentDbSvc)

	t.Run("ReadImageImages", func(t *testing.T) {
		images, err := containerImageQueryRepo.Read()
		if err != nil {
			t.Error(err)
			return
		}

		if len(images) == 0 {
			t.Error("NoImagesFound")
		}
	})
}
