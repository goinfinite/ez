package infra

import (
	"testing"

	testHelpers "github.com/speedianet/sfm/src/devUtils"
)

func TestContainerQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()

	t.Run("GetContainers", func(t *testing.T) {
		containerList, err := ContainerQueryRepo{}.Get()
		if err != nil {
			t.Error(err)
		}

		if len(containerList) == 0 {
			t.Error("NoContainersFound")
		}
	})
}
