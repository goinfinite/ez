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
			valueObject.NewNetworkPortPanic(443),
			valueObject.NewNetworkProtocolPanic("https"),
			[]dto.AddMappingTargetWithoutMappingId{},
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
			nil,
			nil,
		)

		err = mappingCmdRepo.AddTarget(addMappingTarget)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})
}
