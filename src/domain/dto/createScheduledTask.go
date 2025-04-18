package dto

import "github.com/goinfinite/ez/src/domain/valueObject"

type CreateScheduledTask struct {
	Name        valueObject.ScheduledTaskName  `json:"name"`
	Command     valueObject.UnixCommand        `json:"command"`
	Tags        []valueObject.ScheduledTaskTag `json:"tags"`
	TimeoutSecs *valueObject.TimeDuration      `json:"timeoutSecs,omitempty"`
	RunAt       *valueObject.UnixTime          `json:"runAt,omitempty"`
}

func NewCreateScheduledTask(
	name valueObject.ScheduledTaskName,
	command valueObject.UnixCommand,
	tags []valueObject.ScheduledTaskTag,
	timeoutSecs *valueObject.TimeDuration,
	runAt *valueObject.UnixTime,
) CreateScheduledTask {
	return CreateScheduledTask{
		Name:        name,
		Command:     command,
		Tags:        tags,
		TimeoutSecs: timeoutSecs,
		RunAt:       runAt,
	}
}
