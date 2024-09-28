package apiController

import (
	"log/slog"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/domain/valueObject"
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

func (controller *ScheduledTaskController) parseTaskTags(
	rawTags string,
) []valueObject.ScheduledTaskTag {
	taskTags := []valueObject.ScheduledTaskTag{}
	for rawTagIndex, rawTag := range strings.Split(rawTags, ";") {
		taskTag, err := valueObject.NewScheduledTaskTag(rawTag)
		if err != nil {
			slog.Debug("InvalidTaskTag", slog.Int("tagIndex", rawTagIndex))
			continue
		}
		taskTags = append(taskTags, taskTag)
	}

	return taskTags
}

// ReadScheduledTasks	 godoc
// @Summary      ReadScheduledTasks
// @Description  List scheduled tasks.
// @Tags         task
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        taskId query  string  false  "TaskId"
// @Param        taskStatus query  string  false  "TaskStatus"
// @Param        taskTags query  string  false  "TaskTags"
// @Param        startedBeforeAt query  string  false  "StartedBeforeAt"
// @Param        startedAfterAt query  string  false  "StartedAfterAt"
// @Param        finishedBeforeAt query  string  false  "FinishedBeforeAt"
// @Param        finishedAfterAt query  string  false  "FinishedAfterAt"
// @Param        createdBeforeAt query  string  false  "CreatedBeforeAt"
// @Param        createdAfterAt query  string  false  "CreatedAfterAt"
// @Param        pageNumber query  uint  false  "PageNumber (Pagination)"
// @Param        itemsPerPage query  uint  false  "ItemsPerPage (Pagination)"
// @Param        sortBy query  string  false  "SortBy (Pagination)"
// @Param        sortDirection query  string  false  "SortDirection (Pagination)"
// @Param        lastSeenId query  string  false  "LastSeenId (Pagination)"
// @Success      200 {array} dto.ReadScheduledTasksResponse
// @Router       /v1/scheduled-task/ [get]
func (controller *ScheduledTaskController) Read(c echo.Context) error {
	requestBody := map[string]interface{}{}
	queryParameters := []string{
		"taskId", "taskStatus", "taskTags",
		"startedBeforeAt", "startedAfterAt",
		"finishedBeforeAt", "finishedAfterAt",
		"createdBeforeAt", "createdAfterAt",
		"pageNumber", "itemsPerPage", "sortBy", "sortDirection", "lastSeenId",
	}
	for _, paramName := range queryParameters {
		paramValue := c.QueryParam(paramName)
		if paramValue == "" {
			continue
		}

		if paramName != "taskTags" {
			requestBody[paramName] = paramValue
			continue
		}

		requestBody[paramName] = controller.parseTaskTags(paramValue)
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.scheduledTaskService.Read(requestBody),
	)
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
// @Router       /v1/scheduled-task/ [put]
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
