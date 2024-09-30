package infra

import (
	"testing"

	testHelpers "github.com/goinfinite/ez/src/devUtils"
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

func TestContainerProfileCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	containerProfileCmdRepo := NewContainerProfileCmdRepo(persistentDbSvc)
	profileId, _ := valueObject.NewContainerProfileId(2)

	t.Run("CreateContainerProfile", func(t *testing.T) {
		name, _ := valueObject.NewContainerProfileName("testContainerProfile")
		baseMillicores, _ := valueObject.NewMillicores(1000)
		baseMemoryBytes, _ := valueObject.NewByte(1073741824)
		basePerformanceUnits, _ := valueObject.NewStoragePerformanceUnits(1)

		baseSpecs := valueObject.NewContainerSpecs(
			baseMillicores, baseMemoryBytes, basePerformanceUnits,
		)

		maxMillicores, _ := valueObject.NewMillicores(2000)
		maxMemoryBytes, _ := valueObject.NewByte(2147483648)
		maxPerformanceUnits, _ := valueObject.NewStoragePerformanceUnits(2)

		maxSpecs := valueObject.NewContainerSpecs(
			maxMillicores, maxMemoryBytes, maxPerformanceUnits,
		)

		scalingPolicy := valueObject.ScalingPolicy("cpu")
		hostMinCapacityPercent, _ := valueObject.NewHostMinCapacity(10)

		createDto := dto.NewCreateContainerProfile(
			valueObject.SystemAccountId, name, baseSpecs, &maxSpecs,
			&scalingPolicy, nil, nil, nil, &hostMinCapacityPercent,
			valueObject.SystemAccountId, valueObject.SystemIpAddress,
		)

		_, err := containerProfileCmdRepo.Create(createDto)
		if err != nil {
			t.Errorf("CreateContainerProfileFailed: %v", err)
		}
	})

	t.Run("UpdateContainerProfile", func(t *testing.T) {
		maxMillicores, _ := valueObject.NewMillicores(4000)
		maxMemoryBytes, _ := valueObject.NewByte(4294967296)
		maxPerformanceUnits, _ := valueObject.NewStoragePerformanceUnits(3)

		maxSpecs := valueObject.NewContainerSpecs(
			maxMillicores, maxMemoryBytes, maxPerformanceUnits,
		)

		updateDto := dto.NewUpdateContainerProfile(
			valueObject.SystemAccountId, profileId, nil, nil, &maxSpecs,
			nil, nil, nil, nil, nil, valueObject.SystemAccountId,
			valueObject.SystemIpAddress,
		)

		err := containerProfileCmdRepo.Update(updateDto)
		if err != nil {
			t.Errorf("UpdateContainerProfileFailed: %v", err)
		}
	})

	t.Run("DeleteContainerProfile", func(t *testing.T) {
		deleteDto := dto.NewDeleteContainerProfile(
			valueObject.SystemAccountId, profileId,
			valueObject.SystemAccountId, valueObject.SystemIpAddress,
		)
		err := containerProfileCmdRepo.Delete(deleteDto)
		if err != nil {
			t.Errorf("DeleteContainerProfileFailed: %v", err)
		}
	})
}
