package infra

import (
	"errors"
	"testing"

	testHelpers "github.com/speedianet/control/src/devUtils"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

func getScheduledTasks() ([]entity.ScheduledTask, error) {
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	scheduledTaskQueryRepo := NewScheduledTaskQueryRepo(persistentDbSvc)

	scheduledTasks, err := scheduledTaskQueryRepo.Get()
	if err != nil || len(scheduledTasks) == 0 {
		return nil, errors.New("NoScheduledTasksFound")
	}

	return scheduledTasks, nil
}

func TestScheduledTaskQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	scheduledTaskQueryRepo := NewScheduledTaskQueryRepo(persistentDbSvc)

	t.Run("GetScheduledTasks", func(t *testing.T) {
		_, err := scheduledTaskQueryRepo.Get()
		if err != nil {
			t.Error(err)
			return
		}
	})

	t.Run("GetScheduledTaskById", func(t *testing.T) {
		scheduledTasks, err := getScheduledTasks()
		if err != nil {
			t.Error(err)
			return
		}

		_, err = scheduledTaskQueryRepo.GetById(scheduledTasks[0].Id)
		if err != nil {
			t.Error(err)
			return
		}
	})

	t.Run("GetScheduledTasksByStatus", func(t *testing.T) {
		pendingStatus, _ := valueObject.NewScheduledTaskStatus("pending")
		_, err := scheduledTaskQueryRepo.GetByStatus(pendingStatus)
		if err != nil {
			t.Error(err)
			return
		}
	})
}
