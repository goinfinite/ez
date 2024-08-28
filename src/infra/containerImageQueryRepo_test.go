package infra

import (
	"testing"

	testHelpers "github.com/speedianet/control/src/devUtils"
)

func TestContainerImageQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	containerImageQueryRepo := NewContainerImageQueryRepo(persistentDbSvc)

	t.Run("ReadImages", func(t *testing.T) {
		imagesList, err := containerImageQueryRepo.Read()
		if err != nil {
			t.Fatal(err)
		}

		if len(imagesList) == 0 {
			t.Fatal("NoImagesFound")
		}
	})
}
