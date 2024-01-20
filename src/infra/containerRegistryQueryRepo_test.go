package infra

import (
	"testing"

	testHelpers "github.com/speedianet/control/src/devUtils"
	"github.com/speedianet/control/src/domain/valueObject"
)

func TestContainerRegistryQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	dbSvc := testHelpers.GetDbSvc()
	containerQueryRepo := NewContainerRegistryQueryRepo(dbSvc)

	t.Run("GetRegistryImages", func(t *testing.T) {
		imageName, _ := valueObject.NewRegistryImageName("nginx")

		registryImages, err := containerQueryRepo.GetImages(&imageName)
		if err != nil {
			t.Error(err)
		}

		if len(registryImages) == 0 {
			t.Error("NoImagesFound")
		}
	})
}
