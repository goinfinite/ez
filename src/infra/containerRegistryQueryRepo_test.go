package infra

import (
	"testing"

	testHelpers "github.com/speedianet/control/src/devUtils"
	"github.com/speedianet/control/src/domain/valueObject"
)

func TestContainerRegistryQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistDbSvc := testHelpers.GetPersistentDbSvc()
	containerQueryRepo := NewContainerRegistryQueryRepo(persistDbSvc)

	t.Run("GetRegistryImages", func(t *testing.T) {
		imageName, _ := valueObject.NewRegistryImageName("speedianet/os")

		registryImages, err := containerQueryRepo.GetImages(&imageName)
		if err != nil {
			t.Error(err)
		}

		if len(registryImages) == 0 {
			t.Error("NoImagesFound")
		}
	})

	t.Run("GetRegistryTaggedImage", func(t *testing.T) {
		imageAddress, _ := valueObject.NewContainerImageAddress("speedianet/os")

		_, err := containerQueryRepo.GetTaggedImage(imageAddress)
		if err != nil {
			t.Error(err)
		}
	})
}
