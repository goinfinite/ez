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
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	mappingCmdRepo := NewMappingCmdRepo(persistentDbSvc)
	containerQueryRepo := NewContainerQueryRepo(persistentDbSvc)

	t.Run("CreateMapping", func(t *testing.T) {
		accountId, _ := valueObject.NewAccountId(os.Getenv("DUMMY_USER_ID"))
		hostname := valueObject.NewFqdnPanic("speedia.net")

		createMapping := dto.NewCreateMapping(
			accountId,
			&hostname,
			valueObject.NewNetworkPortPanic(80),
			valueObject.NewNetworkProtocolPanic("http"),
			[]valueObject.ContainerId{},
		)

		_, err := mappingCmdRepo.Create(createMapping)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})

	t.Run("CreateTargets", func(t *testing.T) {
		mappingId := valueObject.NewMappingIdPanic(1)
		containers, err := containerQueryRepo.Read()
		if err != nil {
			t.Errorf("ReadContainersFailed: %v", err)
			return
		}

		if len(containers) == 0 {
			t.Errorf("NoContainerRunning")
			return
		}

		containerId := containers[0].Id

		createMappingTarget := dto.NewCreateMappingTarget(
			mappingId,
			containerId,
		)

		err = mappingCmdRepo.CreateTarget(createMappingTarget)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})

	t.Run("DeleteTargets", func(t *testing.T) {
		queryRepo := NewMappingQueryRepo(persistentDbSvc)

		mappings, err := queryRepo.Read()
		if err != nil {
			t.Errorf("ReadMappingsFailed: %v", err)
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

		mappingId := mappings[0].Id
		targetId := mappings[0].Targets[0].Id

		err = mappingCmdRepo.DeleteTarget(mappingId, targetId)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})

	t.Run("DeleteMapping", func(t *testing.T) {
		queryRepo := NewMappingQueryRepo(persistentDbSvc)

		mappings, err := queryRepo.Read()
		if err != nil {
			t.Errorf("ReadMappingsFailed: %v", err)
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
