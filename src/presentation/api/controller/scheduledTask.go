package apiController

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
	apiHelper "github.com/speedianet/control/src/presentation/api/helper"
	"github.com/speedianet/control/src/presentation/service"
)

type ScheduledTaskController struct {
	scheduledTaskService *service.ScheduledTaskService
	persistentDbSvc      *db.PersistentDatabaseService
}

func NewScheduledTaskController(
	persistentDbSvc *db.PersistentDatabaseService,
) *ScheduledTaskController {
	return &ScheduledTaskController{
		scheduledTaskService: service.NewScheduledTaskService(persistentDbSvc),
		persistentDbSvc:      persistentDbSvc,
	}
}

// ReadScheduledTasks	 godoc
// @Summary      ReadScheduledTasks
// @Description  List scheduled tasks.
// @Tags         task
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {array} entity.ScheduledTask
// @Router       /v1/task/ [get]
func (controller *ScheduledTaskController) Read(c echo.Context) error {
	return apiHelper.ServiceResponseWrapper(c, controller.scheduledTaskService.Read())
}

// UpdateScheduledTask godoc
// @Summary      UpdateScheduledTask
// @Description  Reschedule a task or change its status.
// @Tags         task
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        updateScheduledTaskDto 	  body dto.UpdateScheduledTask  true  "UpdateScheduledTask (Only id is required.)"
// @Success      200 {object} object{} "ScheduledTaskUpdated"
// @Router       /v1/task/ [put]
func (controller *ScheduledTaskController) Update(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.scheduledTaskService.Update(requestBody),
	)
}

func (controller *ScheduledTaskController) Run() {
	timer := time.NewTicker(
		time.Duration(int64(useCase.ScheduledTasksRunIntervalSecs)) * time.Second,
	)
	defer timer.Stop()

	scheduledTaskQueryRepo := infra.NewScheduledTaskQueryRepo(controller.persistentDbSvc)
	scheduledTaskCmdRepo := infra.NewScheduledTaskCmdRepo(controller.persistentDbSvc)
	for range timer.C {
		go useCase.RunScheduledTasks(scheduledTaskQueryRepo, scheduledTaskCmdRepo)
	}
}
