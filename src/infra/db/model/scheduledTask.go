package dbModel

import (
	"log/slog"
	"strings"
	"time"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type ScheduledTask struct {
	ID          uint   `gorm:"primarykey"`
	Name        string `gorm:"not null"`
	Status      string `gorm:"not null"`
	Command     string `gorm:"not null"`
	Tags        *string
	TimeoutSecs *uint
	RunAt       *time.Time
	Output      *string
	Error       *string
	StartedAt   *time.Time
	FinishedAt  *time.Time
	ElapsedSecs *uint
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

func (ScheduledTask) TableName() string {
	return "scheduled_tasks"
}

func (ScheduledTask) JoinTags(tags []valueObject.ScheduledTaskTag) string {
	tagsStr := ""
	for _, tag := range tags {
		tagsStr += tag.String() + ";"
	}
	return strings.TrimSuffix(tagsStr, ";")
}

func (ScheduledTask) SplitTags(tagsStr string) []valueObject.ScheduledTaskTag {
	rawTagsList := strings.Split(tagsStr, ";")
	var tags []valueObject.ScheduledTaskTag
	for tagIndex, rawTag := range rawTagsList {
		tag, err := valueObject.NewScheduledTaskTag(rawTag)
		if err != nil {
			slog.Debug(err.Error(), slog.Int("index", tagIndex))
			continue
		}
		tags = append(tags, tag)
	}
	return tags
}

func NewScheduledTask(
	id uint,
	name, status, command string,
	tags []valueObject.ScheduledTaskTag,
	timeoutSecs *uint,
	runAt *time.Time,
	output, err *string,
	startedAt, finishedAt *time.Time,
	elapsedSecs *uint,
) ScheduledTask {
	model := ScheduledTask{
		Name:        name,
		Status:      status,
		Command:     command,
		TimeoutSecs: timeoutSecs,
		RunAt:       runAt,
		Output:      output,
		Error:       err,
		StartedAt:   startedAt,
		FinishedAt:  finishedAt,
		ElapsedSecs: elapsedSecs,
	}

	if id != 0 {
		model.ID = id
	}

	if len(tags) > 0 {
		modelTags := model.JoinTags(tags)
		model.Tags = &modelTags
	}

	return model
}

func (model ScheduledTask) ToEntity() (taskEntity entity.ScheduledTask, err error) {
	id, err := valueObject.NewScheduledTaskId(model.ID)
	if err != nil {
		return taskEntity, err
	}

	name, err := valueObject.NewScheduledTaskName(model.Name)
	if err != nil {
		return taskEntity, err
	}

	status, err := valueObject.NewScheduledTaskStatus(model.Status)
	if err != nil {
		return taskEntity, err
	}

	command, err := valueObject.NewUnixCommand(model.Command)
	if err != nil {
		return taskEntity, err
	}

	tags := []valueObject.ScheduledTaskTag{}
	if model.Tags != nil {
		tags = model.SplitTags(*model.Tags)
	}

	var timeoutSecs *uint
	if model.TimeoutSecs != nil {
		timeoutSecs = model.TimeoutSecs
	}

	var runAtPtr *valueObject.UnixTime
	if model.RunAt != nil {
		runAt := valueObject.NewUnixTimeWithGoTime(*model.RunAt)
		runAtPtr = &runAt
	}

	var outputPtr *valueObject.ScheduledTaskOutput
	if model.Output != nil {
		output, err := valueObject.NewScheduledTaskOutput(*model.Output)
		if err != nil {
			return taskEntity, err
		}
		outputPtr = &output
	}

	var taskErrorPtr *valueObject.ScheduledTaskOutput
	if model.Error != nil {
		taskError, err := valueObject.NewScheduledTaskOutput(*model.Error)
		if err != nil {
			return taskEntity, err
		}
		taskErrorPtr = &taskError
	}

	var startedAtPtr *valueObject.UnixTime
	if model.StartedAt != nil {
		startedAt := valueObject.NewUnixTimeWithGoTime(*model.StartedAt)
		startedAtPtr = &startedAt
	}

	var finishedAtPtr *valueObject.UnixTime
	if model.FinishedAt != nil {
		finishedAt := valueObject.NewUnixTimeWithGoTime(*model.FinishedAt)
		finishedAtPtr = &finishedAt
	}

	createdAt := valueObject.NewUnixTimeWithGoTime(model.CreatedAt)
	updatedAt := valueObject.NewUnixTimeWithGoTime(model.UpdatedAt)

	return entity.NewScheduledTask(
		id, name, status, command, tags, timeoutSecs, runAtPtr, outputPtr,
		taskErrorPtr, startedAtPtr, finishedAtPtr, model.ElapsedSecs, createdAt, updatedAt,
	), nil
}
