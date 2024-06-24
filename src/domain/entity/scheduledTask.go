package entity

import "github.com/speedianet/control/src/domain/valueObject"

type ScheduledTask struct {
	Id          valueObject.ScheduledTaskId      `json:"id"`
	Name        valueObject.ScheduledTaskName    `json:"name"`
	Status      valueObject.ScheduledTaskStatus  `json:"status"`
	Command     valueObject.UnixCommand          `json:"command"`
	Tags        []valueObject.ScheduledTaskTag   `json:"tags"`
	TimeoutSecs *uint                            `json:"timeoutSecs"`
	RunAt       *valueObject.UnixTime            `json:"runAt"`
	Stdout      *valueObject.ScheduledTaskOutput `json:"stdout"`
	Stderr      *valueObject.ScheduledTaskOutput `json:"stderr"`
	CreatedAt   valueObject.UnixTime             `json:"createdAt"`
	UpdatedAt   valueObject.UnixTime             `json:"updatedAt"`
}

func NewScheduledTask(
	id valueObject.ScheduledTaskId,
	name valueObject.ScheduledTaskName,
	status valueObject.ScheduledTaskStatus,
	command valueObject.UnixCommand,
	tags []valueObject.ScheduledTaskTag,
	timeoutSecs *uint,
	runAt *valueObject.UnixTime,
	stdout *valueObject.ScheduledTaskOutput,
	stderr *valueObject.ScheduledTaskOutput,
	createdAt valueObject.UnixTime,
	updatedAt valueObject.UnixTime,
) ScheduledTask {
	return ScheduledTask{
		Id:          id,
		Name:        name,
		Status:      status,
		Command:     command,
		Tags:        tags,
		TimeoutSecs: timeoutSecs,
		RunAt:       runAt,
		Stdout:      stdout,
		Stderr:      stderr,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}
