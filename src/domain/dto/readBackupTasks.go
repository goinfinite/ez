package dto

import (
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type ReadBackupTasksRequest struct {
	Pagination        Pagination                           `json:"pagination"`
	TaskId            *valueObject.BackupTaskId            `json:"taskId"`
	AccountId         *valueObject.AccountId               `json:"accountId"`
	JobId             *valueObject.BackupJobId             `json:"jobId"`
	DestinationId     *valueObject.BackupDestinationId     `json:"destinationId"`
	TaskStatus        *valueObject.BackupTaskStatus        `json:"taskStatus"`
	RetentionStrategy *valueObject.BackupRetentionStrategy `json:"retentionStrategy"`
	ContainerId       *valueObject.ContainerId             `json:"containerId"`
	StartedBeforeAt   *valueObject.UnixTime                `json:"startedBeforeAt"`
	StartedAfterAt    *valueObject.UnixTime                `json:"startedAfterAt"`
	FinishedBeforeAt  *valueObject.UnixTime                `json:"finishedBeforeAt"`
	FinishedAfterAt   *valueObject.UnixTime                `json:"finishedAfterAt"`
	CreatedBeforeAt   *valueObject.UnixTime                `json:"createdBeforeAt"`
	CreatedAfterAt    *valueObject.UnixTime                `json:"createdAfterAt"`
}

type ReadBackupTasksResponse struct {
	Pagination Pagination          `json:"pagination"`
	Tasks      []entity.BackupTask `json:"tasks"`
}
