package infra

import (
	"testing"

	testHelpers "github.com/speedianet/control/src/devUtils"
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/valueObject"
)

func TestContainerProfileCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	containerProfileCmdRepo := NewContainerProfileCmdRepo(persistentDbSvc)

	t.Run("AddContainerProfile", func(t *testing.T) {
		name := valueObject.NewContainerProfileNamePanic("testContainerProfile")

		baseSpecs := valueObject.NewContainerSpecs(
			valueObject.NewCpuCoresCountPanic(1),
			valueObject.NewBytePanic(1073741824),
		)

		maxSpecs := valueObject.NewContainerSpecs(
			valueObject.NewCpuCoresCountPanic(2),
			valueObject.NewBytePanic(2147483648),
		)

		scalingPolicy := valueObject.ScalingPolicy("cpu")
		threshold := uint64(80)
		maxDuration := uint64(3600)
		interval := uint64(86400)
		hostMinCapacityPercent := valueObject.NewHostMinCapacityPanic(10)

		addDto := dto.NewAddContainerProfile(
			name,
			baseSpecs,
			&maxSpecs,
			&scalingPolicy,
			&threshold,
			&maxDuration,
			&interval,
			&hostMinCapacityPercent,
		)

		err := containerProfileCmdRepo.Add(addDto)
		if err != nil {
			t.Errorf("AddContainerProfileFailed: %v", err)
		}
	})

	t.Run("UpdateContainerProfile", func(t *testing.T) {
		id := valueObject.NewContainerProfileIdPanic(2)

		maxSpecs := valueObject.NewContainerSpecs(
			valueObject.NewCpuCoresCountPanic(4),
			valueObject.NewBytePanic(4294967296),
		)

		addDto := dto.NewUpdateContainerProfile(
			id,
			nil,
			nil,
			&maxSpecs,
			nil,
			nil,
			nil,
			nil,
			nil,
		)

		err := containerProfileCmdRepo.Update(addDto)
		if err != nil {
			t.Errorf("UpdateContainerProfileFailed: %v", err)
		}
	})

	t.Run("DeleteContainerProfile", func(t *testing.T) {
		id := valueObject.NewContainerProfileIdPanic(2)

		err := containerProfileCmdRepo.Delete(id)
		if err != nil {
			t.Errorf("DeleteContainerProfileFailed: %v", err)
		}
	})
}
