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
	Output      *valueObject.ScheduledTaskOutput `json:"output"`
	Error       *valueObject.ScheduledTaskOutput `json:"err"`
	StartedAt   *valueObject.UnixTime            `json:"startedAt"`
	FinishedAt  *valueObject.UnixTime            `json:"finishedAt"`
	ElapsedSecs *uint                            `json:"elapsedSecs"`
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
	output *valueObject.ScheduledTaskOutput,
	err *valueObject.ScheduledTaskOutput,
	startedAt *valueObject.UnixTime,
	finishedAt *valueObject.UnixTime,
	elapsedSecs *uint,
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
		Output:      output,
		Error:       err,
		StartedAt:   startedAt,
		FinishedAt:  finishedAt,
		ElapsedSecs: elapsedSecs,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}
