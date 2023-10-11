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

		addDto := dto.AddResourceProfile{
			Name:                   "testResourceProfile",
			BaseSpecs:              baseSpecs,
			MaxSpecs:               &maxSpecs,
			ScalingPolicy:          &scalingPolicy,
			ScalingThreshold:       &threshold,
			ScalingMaxDurationSecs: &maxDuration,
			ScalingIntervalSecs:    &interval,
			HostMinCapacityPercent: &hostMinCapacityPercent,
		}

		err := resourceProfileCmdRepo.Add(addDto)
		if err != nil {
			t.Errorf("AddResourceProfileFailed: %v", err)
		}
	})
}
