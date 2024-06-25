package infra

import (
	"testing"

	testHelpers "github.com/speedianet/control/src/devUtils"
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/valueObject"
)

func TestScheduledTaskCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	scheduledTaskCmdRepo := NewScheduledTaskCmdRepo(persistentDbSvc)
	scheduledTaskQueryRepo := NewScheduledTaskQueryRepo(persistentDbSvc)

	t.Run("CreateScheduledTask", func(t *testing.T) {
		name, _ := valueObject.NewScheduledTaskName("test")
		command, _ := valueObject.NewUnixCommand("echo test")
		containerTag, _ := valueObject.NewScheduledTaskTag("container")
		tags := []valueObject.ScheduledTaskTag{containerTag}
		timeoutSecs := uint(60)
		runAt := valueObject.NewUnixTimeNow()

		createDto := dto.NewCreateScheduledTask(
			name, command, tags, &timeoutSecs, &runAt,
		)

		err := scheduledTaskCmdRepo.Create(createDto)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})

	t.Run("UpdateScheduledTask", func(t *testing.T) {
		taskEntities, err := scheduledTaskQueryRepo.Get()
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
			return
		}

		if len(taskEntities) == 0 {
			t.Error("NoScheduledTasksFound")
			return
		}

		newStatus, _ := valueObject.NewScheduledTaskStatus("pending")
		updateDto := dto.NewUpdateScheduledTask(
			taskEntities[0].Id, &newStatus, nil,
		)

		err = scheduledTaskCmdRepo.Update(updateDto)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})
}
