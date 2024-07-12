package infra

import (
	"strconv"
	"time"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
	dbModel "github.com/speedianet/control/src/infra/db/model"
	infraHelper "github.com/speedianet/control/src/infra/helper"
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

	var runAtPtr *time.Time
	if createDto.RunAt != nil {
		runAt := time.Unix(createDto.RunAt.Read(), 0)
		runAtPtr = &runAt
	}

	scheduledTaskModel := dbModel.NewScheduledTask(
		0, createDto.Name.String(), newTaskStatus.String(), createDto.Command.String(),
		createDto.Tags, createDto.TimeoutSecs, runAtPtr, nil, nil,
	)

	return repo.persistentDbSvc.Handler.Create(&scheduledTaskModel).Error
}

func (repo *ScheduledTaskCmdRepo) Update(
	updateDto dto.UpdateScheduledTask,
) error {
	updateMap := map[string]interface{}{}

	if updateDto.Status != nil {
		updateMap["status"] = updateDto.Status.String()
	}

	if updateDto.RunAt != nil {
		updateMap["run_at"] = updateDto.RunAt.GetAsGoTime()
	}

	if len(updateMap) == 0 {
		return nil
	}

	return repo.persistentDbSvc.Handler.
		Model(&dbModel.ScheduledTask{}).
		Where("id = ?", updateDto.Id).
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

	cmdWithTimeout := "timeout --kill-after=10s " + timeoutStr + " " + pendingTask.Command.String()
	rawOutput, rawError := infraHelper.RunCmdWithSubShell(cmdWithTimeout)

	finalStatus, _ := valueObject.NewScheduledTaskStatus("completed")
	if rawError != nil {
		finalStatus, _ = valueObject.NewScheduledTaskStatus("failed")
	}

	updateMap := map[string]interface{}{
		"status": finalStatus.String(),
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
		Where("id = ?", pendingTask.Id).
		Updates(updateMap).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo *ScheduledTaskCmdRepo) Delete(id valueObject.ScheduledTaskId) error {
	return repo.persistentDbSvc.Handler.
		Where("id = ?", id).
		Delete(&dbModel.ScheduledTask{}).Error
}
