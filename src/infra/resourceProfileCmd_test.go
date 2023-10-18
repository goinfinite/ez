package infra

import (
	"testing"

	testHelpers "github.com/speedianet/sfm/src/devUtils"
	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

func TestResourceProfileCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	dbSvc := testHelpers.GetDbSvc()
	resourceProfileCmdRepo := NewResourceProfileCmdRepo(dbSvc)

	t.Run("AddResourceProfile", func(t *testing.T) {
		name := valueObject.NewResourceProfileNamePanic("testResourceProfile")

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

		addDto := dto.NewAddResourceProfile(
			name,
			baseSpecs,
			&maxSpecs,
			&scalingPolicy,
			&threshold,
			&maxDuration,
			&interval,
			&hostMinCapacityPercent,
		)

		err := resourceProfileCmdRepo.Add(addDto)
		if err != nil {
			t.Errorf("AddResourceProfileFailed: %v", err)
		}
	})

	t.Run("UpdateResourceProfile", func(t *testing.T) {
		id := valueObject.NewResourceProfileIdPanic(2)

		maxSpecs := valueObject.NewContainerSpecs(
			valueObject.NewCpuCoresCountPanic(4),
			valueObject.NewBytePanic(4294967296),
		)

		addDto := dto.NewUpdateResourceProfile(
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

		err := resourceProfileCmdRepo.Update(addDto)
		if err != nil {
			t.Errorf("UpdateResourceProfileFailed: %v", err)
		}
	})

	t.Run("DeleteResourceProfile", func(t *testing.T) {
		id := valueObject.NewResourceProfileIdPanic(2)

		err := resourceProfileCmdRepo.Delete(id)
		if err != nil {
			t.Errorf("DeleteResourceProfileFailed: %v", err)
		}
	})
}
