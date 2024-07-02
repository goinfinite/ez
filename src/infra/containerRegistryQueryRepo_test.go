package infra

import (
	"testing"

	testHelpers "github.com/speedianet/control/src/devUtils"
	"github.com/speedianet/control/src/domain/valueObject"
)

func TestContainerRegistryQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	containerQueryRepo := NewContainerRegistryQueryRepo(persistentDbSvc)

	t.Run("ReadRegistryImages", func(t *testing.T) {
		imageName, _ := valueObject.NewRegistryImageName("speedianet/os")

		registryImages, err := containerQueryRepo.ReadImages(&imageName)
		if err != nil {
			t.Error(err)
		}

		if len(registryImages) == 0 {
			t.Error("NoImagesFound")
		}
	})

	t.Run("ReadRegistryTaggedImage", func(t *testing.T) {
		imageAddress, _ := valueObject.NewContainerImageAddress("speedianet/os")

		_, err := containerQueryRepo.ReadTaggedImage(imageAddress)
		if err != nil {
			t.Error(err)
		}
	})
}
