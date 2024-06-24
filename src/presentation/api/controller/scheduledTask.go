package apiController

import (
	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/infra/db"
	apiHelper "github.com/speedianet/control/src/presentation/api/helper"
	"github.com/speedianet/control/src/presentation/service"
)

type ScheduledTaskController struct {
	scheduledTaskService *service.ScheduledTaskService
}

func NewScheduledTaskController(
	persistentDbSvc *db.PersistentDatabaseService,
) *ScheduledTaskController {
	return &ScheduledTaskController{
		scheduledTaskService: service.NewScheduledTaskService(persistentDbSvc),
	}
}

// ReadScheduledTasks	 godoc
// @Summary      ReadScheduledTasks
// @Description  List scheduledTasks.
// @Tags         task
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {array} entity.ScheduledTask
// @Router       /v1/task/ [get]
func (controller *ScheduledTaskController) Read(c echo.Context) error {
	return apiHelper.ServiceResponseWrapper(c, controller.scheduledTaskService.Read())
}
