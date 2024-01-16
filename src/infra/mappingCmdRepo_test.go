package infra

import (
	"os"
	"testing"

	testHelpers "github.com/speedianet/control/src/devUtils"
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/valueObject"
)

func TestMappingCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	dbSvc := testHelpers.GetDbSvc()
	mappingCmdRepo := NewMappingCmdRepo(dbSvc)
	containerQueryRepo := NewContainerQueryRepo(dbSvc)

	t.Run("AddMapping", func(t *testing.T) {
		accountId := valueObject.NewAccountIdPanic(os.Getenv("DUMMY_USER_ID"))
		hostname := valueObject.NewFqdnPanic("speedia.net")

		addMapping := dto.NewAddMapping(
			accountId,
			&hostname,
			valueObject.NewNetworkPortPanic(80),
			valueObject.NewNetworkProtocolPanic("http"),
			[]valueObject.ContainerId{},
		)

		_, err := mappingCmdRepo.Add(addMapping)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})

	t.Run("AddTargets", func(t *testing.T) {
		mappingId := valueObject.NewMappingIdPanic(1)
		containers, err := containerQueryRepo.Get()
		if err != nil {
			t.Errorf("GetContainersFailed: %v", err)
			return
		}

		if len(containers) == 0 {
			t.Errorf("NoContainerRunning")
			return
		}

		containerId := containers[0].Id

		addMappingTarget := dto.NewAddMappingTarget(
			mappingId,
			containerId,
		)

		err = mappingCmdRepo.AddTarget(addMappingTarget)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})

	t.Run("DeleteTargets", func(t *testing.T) {
		queryRepo := NewMappingQueryRepo(dbSvc)

		mappings, err := queryRepo.Get()
		if err != nil {
			t.Errorf("GetMappingsFailed: %v", err)
			return
		}

		if len(mappings) == 0 {
			t.Errorf("NoMappingFound")
			return
		}

		if len(mappings[0].Targets) == 0 {
			t.Errorf("NoTargetFound")
			return
		}

		targetId := mappings[0].Targets[0].Id

		err = mappingCmdRepo.DeleteTarget(targetId)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})

	t.Run("DeleteMapping", func(t *testing.T) {
		queryRepo := NewMappingQueryRepo(dbSvc)

		mappings, err := queryRepo.Get()
		if err != nil {
			t.Errorf("GetMappingsFailed: %v", err)
			return
		}

		if len(mappings) == 0 {
			t.Errorf("NoMappingFound")
			return
		}

		mappingId := mappings[0].Id

		err = mappingCmdRepo.Delete(mappingId)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})
}
