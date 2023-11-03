package infra

import (
	"testing"

	testHelpers "github.com/goinfinite/fleet/src/devUtils"
)

func TestContainerQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	dbSvc := testHelpers.GetDbSvc()
	containerQueryRepo := NewContainerQueryRepo(dbSvc)

	t.Run("GetContainers", func(t *testing.T) {
		containerList, err := containerQueryRepo.Get()
		if err != nil {
			t.Error(err)
		}

		if len(containerList) == 0 {
			t.Error("NoContainersFound")
		}
	})

	t.Run("GetContainersWithUsage", func(t *testing.T) {
		containerList, err := containerQueryRepo.GetWithUsage()
		if err != nil {
			t.Error(err)
		}

		if len(containerList) == 0 {
			t.Error("NoContainersFound")
		}
	})
}
