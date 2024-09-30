package infra

import (
	"strconv"
	"time"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
	"github.com/goinfinite/ez/src/infra/db"
	dbModel "github.com/goinfinite/ez/src/infra/db/model"
	infraHelper "github.com/goinfinite/ez/src/infra/helper"
)

type ScheduledTaskCmdRepo struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewScheduledTaskCmdRepo(
	persistentDbSvc *db.PersistentDatabaseService,
) *ScheduledTaskCmdRepo {
	return &ScheduledTaskCmdRepo{persistentDbSvc: persistentDbSvc}
}

func (repo *ScheduledTaskCmdRepo) Create(
	createDto dto.CreateScheduledTask,
) error {
	newTaskStatus, _ := valueObject.NewScheduledTaskStatus("pending")

	taskTagsModels := []dbModel.ScheduledTaskTag{}
	for _, taskTag := range createDto.Tags {
		taskTagModel := dbModel.ScheduledTaskTag{
			Tag: taskTag.String(),
		}
		taskTagsModels = append(taskTagsModels, taskTagModel)
	}

	var runAtPtr *time.Time
	if createDto.RunAt != nil {
		runAt := time.Unix(createDto.RunAt.Read(), 0)
		runAtPtr = &runAt
	}

	scheduledTaskModel := dbModel.NewScheduledTask(
		0, createDto.Name.String(), newTaskStatus.String(), createDto.Command.String(),
		taskTagsModels, createDto.TimeoutSecs, runAtPtr, nil, nil, nil, nil, nil,
	)

	return repo.persistentDbSvc.Handler.Create(&scheduledTaskModel).Error
}

func (repo *ScheduledTaskCmdRepo) Update(
	updateDto dto.UpdateScheduledTask,
) error {
	updateMap := map[string]interface{}{}

	if updateDto.Status != nil {
		updateMap["status"] = updateDto.Status.String()
		updateMap["run_at"] = nil
		updateMap["output"] = nil
		updateMap["error"] = nil
		updateMap["started_at"] = nil
		updateMap["finished_at"] = nil
		updateMap["elapsed_secs"] = nil
	}

	if updateDto.RunAt != nil {
		updateMap["run_at"] = updateDto.RunAt.GetAsGoTime()
	}

	if len(updateMap) == 0 {
		return nil
	}

	return repo.persistentDbSvc.Handler.
		Model(&dbModel.ScheduledTask{}).
		Where("id = ?", updateDto.TaskId).
		Updates(updateMap).Error
}

func (repo *ScheduledTaskCmdRepo) Run(
	pendingTask entity.ScheduledTask,
) error {
	runningStatus, _ := valueObject.NewScheduledTaskStatus("running")
	updateDto := dto.NewUpdateScheduledTask(pendingTask.Id, &runningStatus, nil)
	err := repo.Update(updateDto)
	if err != nil {
		return err
	}

	timeoutStr := "300"
	if pendingTask.TimeoutSecs != nil {
		timeoutStr = strconv.FormatUint(uint64(*pendingTask.TimeoutSecs), 10)
	}

	startedAtUnixTime := valueObject.NewUnixTimeNow()

	cmdWithTimeout := "timeout --kill-after=10s " + timeoutStr + " " + pendingTask.Command.String()
	rawOutput, rawError := infraHelper.RunCmdWithSubShell(cmdWithTimeout)

	finalStatus, _ := valueObject.NewScheduledTaskStatus("completed")
	if rawError != nil {
		finalStatus, _ = valueObject.NewScheduledTaskStatus("failed")
	}

	finishedAtUnixTime := valueObject.NewUnixTimeNow()
	elapsedSecs := uint(finishedAtUnixTime.Read() - startedAtUnixTime.Read())

	updateMap := map[string]interface{}{
		"status":       finalStatus.String(),
		"started_at":   startedAtUnixTime.GetAsGoTime(),
		"finished_at":  finishedAtUnixTime.GetAsGoTime(),
		"elapsed_secs": elapsedSecs,
	}

	if len(rawOutput) > 0 {
		taskOutput, err := valueObject.NewScheduledTaskOutput(rawOutput)
		if err != nil {
			return err
		}
		updateMap["output"] = taskOutput.String()
	}

	if rawError != nil {
		taskError, err := valueObject.NewScheduledTaskOutput(rawError.Error())
		if err != nil {
			return err
		}
		updateMap["error"] = taskError.String()
	}

	err = repo.persistentDbSvc.Handler.
		Model(&dbModel.ScheduledTask{}).
		Preload("Tags").
		Where("id = ?", pendingTask.Id).
		Updates(updateMap).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo *ScheduledTaskCmdRepo) Delete(id valueObject.ScheduledTaskId) error {
	return repo.persistentDbSvc.Handler.
		Model(&dbModel.ScheduledTask{}).
		Delete("id = ?", id.Uint64()).Error
}
