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

		port, _ := valueObject.NewNetworkPort(80)
		protocol, _ := valueObject.NewNetworkProtocol("http")

		createMapping := dto.NewCreateMapping(
			accountId, &hostname, port, protocol, []valueObject.ContainerId{},
			valueObject.SystemAccountId, valueObject.SystemIpAddress,
		)

		_, err := mappingCmdRepo.Create(createMapping)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})

	t.Run("CreateTargets", func(t *testing.T) {
		mappingId, _ := valueObject.NewMappingId(1)
		containers, err := containerQueryRepo.Read()
		if err != nil {
			t.Errorf("ReadContainersFailed: %v", err)
			return
		}

		if len(containers) == 0 {
			t.Errorf("NoContainerRunning")
			return
		}

		accountId := containers[0].AccountId
		containerId := containers[0].Id

		createMappingTarget := dto.NewCreateMappingTarget(
			accountId, mappingId, containerId,
			valueObject.SystemAccountId, valueObject.SystemIpAddress,
		)

		_, err = mappingCmdRepo.CreateTarget(createMappingTarget)
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

		accountId := mappings[0].AccountId
		mappingId := mappings[0].Id
		targetId := mappings[0].Targets[0].Id
		deleteDto := dto.NewDeleteMappingTarget(
			accountId, mappingId, targetId,
			valueObject.SystemAccountId, valueObject.SystemIpAddress,
		)

		err = mappingCmdRepo.DeleteTarget(deleteDto)
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

		deleteDto := dto.DeleteMapping{MappingId: mappings[0].Id}
		err = mappingCmdRepo.Delete(deleteDto)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})
}
