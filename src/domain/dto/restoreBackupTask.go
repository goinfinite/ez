package dto

import (
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type RestoreBackupTaskRequest struct {
	TaskId                          *valueObject.BackupTaskId        `json:"taskId,omitempty"`
	ArchiveId                       *valueObject.BackupTaskArchiveId `json:"archiveId,omitempty"`
	ShouldReplaceExistingContainers *bool                            `json:"shouldReplaceExistingContainers,omitempty"`
	ShouldRestoreMappings           *bool                            `json:"shouldRestoreMappings,omitempty"`
	TimeoutSecs                     *uint32                          `json:"timeoutSecs,omitempty"`
	ContainerAccountIds             []valueObject.AccountId          `json:"containerAccountIds"`
	ContainerIds                    []valueObject.ContainerId        `json:"containerIds"`
	ExceptContainerAccountIds       []valueObject.AccountId          `json:"exceptContainerAccountIds"`
	ExceptContainerIds              []valueObject.ContainerId        `json:"exceptContainerIds"`
	OperatorAccountId               valueObject.AccountId            `json:"-"`
	OperatorIpAddress               valueObject.IpAddress            `json:"-"`
}

type RestoreBackupTaskResponse struct {
	SuccessfulContainerIds  []valueObject.ContainerId      `json:"successfulContainerIds"`
	FailedContainerImageIds []valueObject.ContainerImageId `json:"failedContainerImageIds"`
}
