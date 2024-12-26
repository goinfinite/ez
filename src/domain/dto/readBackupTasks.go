package dto

import (
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type ReadBackupTasksRequest struct {
	Pagination        Pagination                           `json:"pagination"`
	TaskId            *valueObject.BackupTaskId            `json:"taskId,omitempty"`
	AccountId         *valueObject.AccountId               `json:"accountId,omitempty"`
	JobId             *valueObject.BackupJobId             `json:"jobId,omitempty"`
	DestinationId     *valueObject.BackupDestinationId     `json:"destinationId,omitempty"`
	TaskStatus        *valueObject.BackupTaskStatus        `json:"taskStatus,omitempty"`
	RetentionStrategy *valueObject.BackupRetentionStrategy `json:"retentionStrategy,omitempty"`
	ContainerId       *valueObject.ContainerId             `json:"containerId,omitempty"`
	StartedBeforeAt   *valueObject.UnixTime                `json:"startedBeforeAt,omitempty"`
	StartedAfterAt    *valueObject.UnixTime                `json:"startedAfterAt,omitempty"`
	FinishedBeforeAt  *valueObject.UnixTime                `json:"finishedBeforeAt,omitempty"`
	FinishedAfterAt   *valueObject.UnixTime                `json:"finishedAfterAt,omitempty"`
	CreatedBeforeAt   *valueObject.UnixTime                `json:"createdBeforeAt,omitempty"`
	CreatedAfterAt    *valueObject.UnixTime                `json:"createdAfterAt,omitempty"`
}

type ReadBackupTasksResponse struct {
	Pagination Pagination          `json:"pagination"`
	Tasks      []entity.BackupTask `json:"tasks"`
}
