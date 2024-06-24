package cliController

import (
	"github.com/speedianet/control/src/infra/db"
	cliHelper "github.com/speedianet/control/src/presentation/cli/helper"
	"github.com/speedianet/control/src/presentation/service"
	"github.com/spf13/cobra"
)

type ScheduledTaskController struct {
	persistentDbSvc      *db.PersistentDatabaseService
	scheduledTaskService *service.ScheduledTaskService
}

func NewScheduledTaskController(
	persistentDbSvc *db.PersistentDatabaseService,
) *ScheduledTaskController {
	return &ScheduledTaskController{
		persistentDbSvc:      persistentDbSvc,
		scheduledTaskService: service.NewScheduledTaskService(persistentDbSvc),
	}
}

func (controller *ScheduledTaskController) Read() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "ReadScheduledTasks",
		Run: func(cmd *cobra.Command, args []string) {
			cliHelper.ServiceResponseWrapper(controller.scheduledTaskService.Read())
		},
	}

	return cmd
}
