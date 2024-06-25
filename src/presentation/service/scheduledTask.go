package service

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
	serviceHelper "github.com/speedianet/control/src/presentation/service/helper"
)

type ScheduledTaskService struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewScheduledTaskService(
	persistentDbSvc *db.PersistentDatabaseService,
) *ScheduledTaskService {
	return &ScheduledTaskService{
		persistentDbSvc: persistentDbSvc,
	}
}

func (service *ScheduledTaskService) Read() ServiceOutput {
	scheduledTaskQueryRepo := infra.NewScheduledTaskQueryRepo(service.persistentDbSvc)
	scheduledTasksList, err := useCase.GetScheduledTasks(scheduledTaskQueryRepo)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, scheduledTasksList)
}

func (service *ScheduledTaskService) Update(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"id"}

	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	taskId, err := valueObject.NewScheduledTaskId(input["id"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var taskStatusPtr *valueObject.ScheduledTaskStatus
	if _, exists := input["status"]; exists {
		taskStatus, err := valueObject.NewScheduledTaskStatus(input["status"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		taskStatusPtr = &taskStatus
	}

	var runAtPtr *valueObject.UnixTime
	if _, exists := input["runAt"]; exists {
		runAt, err := valueObject.NewUnixTime(input["runAt"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		runAtPtr = &runAt
	}

	updateDto := dto.NewUpdateScheduledTask(
		taskId,
		taskStatusPtr,
		runAtPtr,
	)

	scheduledTaskQueryRepo := infra.NewScheduledTaskQueryRepo(service.persistentDbSvc)
	scheduledTaskCmdRepo := infra.NewScheduledTaskCmdRepo(service.persistentDbSvc)

	err = useCase.UpdateScheduledTask(
		scheduledTaskQueryRepo,
		scheduledTaskCmdRepo,
		updateDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "ScheduledTaskUpdated")
}
