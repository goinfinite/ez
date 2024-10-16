package infra

import (
	"testing"

	testHelpers "github.com/goinfinite/ez/src/devUtils"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

func TestContainerRegistryQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	containerQueryRepo := NewContainerRegistryQueryRepo(persistentDbSvc)

	t.Run("ReadRegistryImages", func(t *testing.T) {
		imageName, _ := valueObject.NewRegistryImageName("goinfinite/os")

		registryImages, err := containerQueryRepo.ReadImages(&imageName)
		if err != nil {
			t.Error(err)
		}

		if len(registryImages) == 0 {
			t.Error("NoImagesFound")
		}
	})

	t.Run("ReadRegistryTaggedImage", func(t *testing.T) {
		imageAddress, _ := valueObject.NewContainerImageAddress("goinfinite/os")

		_, err := containerQueryRepo.ReadTaggedImage(imageAddress)
		if err != nil {
			t.Error(err)
		}
	})
}
