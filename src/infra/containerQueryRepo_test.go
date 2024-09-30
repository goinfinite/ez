package infra

import (
	"testing"

	testHelpers "github.com/goinfinite/ez/src/devUtils"
)

func TestContainerQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	containerQueryRepo := NewContainerQueryRepo(persistentDbSvc)

	t.Run("ReadContainers", func(t *testing.T) {
		containerList, err := containerQueryRepo.Read()
		if err != nil {
			t.Error(err)
		}

		if len(containerList) == 0 {
			t.Error("NoContainersFound")
		}
	})

	t.Run("ReadContainersWithMetrics", func(t *testing.T) {
		containerList, err := containerQueryRepo.ReadWithMetrics()
		if err != nil {
			t.Error(err)
		}

		if len(containerList) == 0 {
			t.Error("NoContainersFound")
		}
	})
}
