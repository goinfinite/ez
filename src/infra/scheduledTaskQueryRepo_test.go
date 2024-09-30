package infra

import (
	"errors"
	"testing"

	testHelpers "github.com/goinfinite/ez/src/devUtils"
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/useCase"
)

func readScheduledTasks() ([]entity.ScheduledTask, error) {
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	scheduledTaskQueryRepo := NewScheduledTaskQueryRepo(persistentDbSvc)

	readDto := dto.ReadScheduledTasksRequest{
		Pagination: useCase.ScheduledTasksDefaultPagination,
	}

	responseDto, err := scheduledTaskQueryRepo.Read(readDto)
	if err != nil {
		return nil, err
	}

	if len(responseDto.Tasks) == 0 {
		return nil, errors.New("NoScheduledTasksFound")
	}

	return responseDto.Tasks, nil
}

func TestScheduledTaskQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()

	t.Run("ReadScheduledTasks", func(t *testing.T) {
		_, err := readScheduledTasks()
		if err != nil {
			t.Error(err)
		}
	})
}
