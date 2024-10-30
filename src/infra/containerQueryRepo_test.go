package infra

import (
	"testing"

	testHelpers "github.com/goinfinite/ez/src/devUtils"
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/useCase"
)

func TestContainerQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	containerQueryRepo := NewContainerQueryRepo(persistentDbSvc)

	t.Run("ReadContainers", func(t *testing.T) {
		requestDto := dto.ReadContainersRequest{
			Pagination: useCase.ContainersDefaultPagination,
		}

		responseDto, err := containerQueryRepo.Read(requestDto)
		if err != nil {
			t.Error(err)
		}

		if len(responseDto.Containers) == 0 {
			t.Error("NoContainersFound")
		}
	})

	t.Run("ReadContainersWithMetrics", func(t *testing.T) {
		withMetrics := true
		requestDto := dto.ReadContainersRequest{
			Pagination:  useCase.ContainersDefaultPagination,
			WithMetrics: &withMetrics,
		}

		responseDto, err := containerQueryRepo.Read(requestDto)
		if err != nil {
			t.Error(err)
		}

		if len(responseDto.ContainersWithMetrics) == 0 {
			t.Error("NoContainersFound")
		}
	})
}
