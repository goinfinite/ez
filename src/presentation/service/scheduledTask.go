package service

import (
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/useCase"
	"github.com/goinfinite/ez/src/domain/valueObject"
	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
	"github.com/goinfinite/ez/src/infra"
	"github.com/goinfinite/ez/src/infra/db"
	serviceHelper "github.com/goinfinite/ez/src/presentation/service/helper"
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

func (service *ScheduledTaskService) Read(input map[string]interface{}) ServiceOutput {
	var taskIdPtr *valueObject.ScheduledTaskId
	if input["id"] != nil {
		input["taskId"] = input["id"]
	}
	if input["taskId"] != nil {
		taskId, err := valueObject.NewScheduledTaskId(input["taskId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		taskIdPtr = &taskId
	}

	var taskNamePtr *valueObject.ScheduledTaskName
	if input["taskName"] != nil {
		taskName, err := valueObject.NewScheduledTaskName(input["taskName"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		taskNamePtr = &taskName
	}

	var taskStatusPtr *valueObject.ScheduledTaskStatus
	if input["taskStatus"] != nil {
		taskStatus, err := valueObject.NewScheduledTaskStatus(input["taskStatus"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		taskStatusPtr = &taskStatus
	}

	taskTags := []valueObject.ScheduledTaskTag{}
	if input["taskTags"] != nil {
		var assertOk bool
		taskTags, assertOk = input["taskTags"].([]valueObject.ScheduledTaskTag)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidTaskTags")
		}
	}

	timeParamNames := []string{
		"startedBeforeAt", "startedAfterAt",
		"finishedBeforeAt", "finishedAfterAt",
		"createdBeforeAt", "createdAfterAt",
	}
	timeParamPtrs := serviceHelper.TimeParamsParser(timeParamNames, input)

	paginationDto := useCase.ScheduledTasksDefaultPagination
	if input["pageNumber"] != nil {
		pageNumber, err := voHelper.InterfaceToUint32(input["pageNumber"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidPageNumber")
		}
		paginationDto.PageNumber = pageNumber
	}

	if input["itemsPerPage"] != nil {
		itemsPerPage, err := voHelper.InterfaceToUint16(input["itemsPerPage"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidItemsPerPage")
		}
		paginationDto.ItemsPerPage = itemsPerPage
	}

	if input["sortBy"] != nil {
		sortBy, err := valueObject.NewPaginationSortBy(input["sortBy"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		paginationDto.SortBy = &sortBy
	}

	if input["sortDirection"] != nil {
		sortDirection, err := valueObject.NewPaginationSortDirection(input["sortDirection"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		paginationDto.SortDirection = &sortDirection
	}

	if input["lastSeenId"] != nil {
		lastSeenId, err := valueObject.NewPaginationLastSeenId(input["lastSeenId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		paginationDto.LastSeenId = &lastSeenId
	}

	readDto := dto.ReadScheduledTasksRequest{
		Pagination:       paginationDto,
		TaskId:           taskIdPtr,
		TaskName:         taskNamePtr,
		TaskStatus:       taskStatusPtr,
		TaskTags:         taskTags,
		StartedBeforeAt:  timeParamPtrs["startedBeforeAt"],
		StartedAfterAt:   timeParamPtrs["startedAfterAt"],
		FinishedBeforeAt: timeParamPtrs["finishedBeforeAt"],
		FinishedAfterAt:  timeParamPtrs["finishedAfterAt"],
		CreatedBeforeAt:  timeParamPtrs["createdBeforeAt"],
		CreatedAfterAt:   timeParamPtrs["createdAfterAt"],
	}

	scheduledTaskQueryRepo := infra.NewScheduledTaskQueryRepo(service.persistentDbSvc)
	scheduledTasksList, err := useCase.ReadScheduledTasks(scheduledTaskQueryRepo, readDto)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, scheduledTasksList)
}

func (service *ScheduledTaskService) Update(input map[string]interface{}) ServiceOutput {
	if input["id"] != nil {
		input["taskId"] = input["id"]
	}

	requiredParams := []string{"taskId"}

	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	taskId, err := valueObject.NewScheduledTaskId(input["taskId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var taskStatusPtr *valueObject.ScheduledTaskStatus
	if input["status"] != nil {
		taskStatus, err := valueObject.NewScheduledTaskStatus(input["status"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		taskStatusPtr = &taskStatus
	}

	var runAtPtr *valueObject.UnixTime
	if input["runAt"] != nil {
		runAt, err := valueObject.NewUnixTime(input["runAt"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		runAtPtr = &runAt
	}

	updateDto := dto.NewUpdateScheduledTask(
		taskId, taskStatusPtr, runAtPtr,
	)

	scheduledTaskQueryRepo := infra.NewScheduledTaskQueryRepo(service.persistentDbSvc)
	scheduledTaskCmdRepo := infra.NewScheduledTaskCmdRepo(service.persistentDbSvc)

	err = useCase.UpdateScheduledTask(
		scheduledTaskQueryRepo, scheduledTaskCmdRepo, updateDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "ScheduledTaskUpdated")
}
