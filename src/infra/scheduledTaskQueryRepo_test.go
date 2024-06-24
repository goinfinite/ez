package infra

import (
	"testing"

	testHelpers "github.com/speedianet/control/src/devUtils"
)

func TestScheduledTaskQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	scheduledTaskQueryRepo := NewScheduledTaskQueryRepo(persistentDbSvc)

	t.Run("GetScheduledTasks", func(t *testing.T) {
		scheduledTaskList, err := scheduledTaskQueryRepo.Get()
		if err != nil {
			t.Error(err)
			return
		}

		if len(scheduledTaskList) == 0 {
			t.Error("NoScheduledTasksFound")
			return
		}
	})
}
