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

	t.Run("CreateScheduledTask", func(t *testing.T) {
		name, _ := valueObject.NewScheduledTaskName("test")
		command, _ := valueObject.NewUnixCommand("echo test")
		containerTag, _ := valueObject.NewScheduledTaskTag("container")
		tags := []valueObject.ScheduledTaskTag{containerTag}
		timeoutSecs := uint(60)
		runAt := valueObject.NewUnixTimeNow()

		createScheduledTask := dto.NewCreateScheduledTask(
			name, command, tags, &timeoutSecs, &runAt,
		)

		err := scheduledTaskCmdRepo.Create(createScheduledTask)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})
}
