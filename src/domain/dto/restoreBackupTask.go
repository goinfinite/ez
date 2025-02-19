package dto

import (
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type RestoreBackupTask struct {
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

func NewRestoreBackupTask(
	taskId *valueObject.BackupTaskId,
	archiveId *valueObject.BackupTaskArchiveId,
	shouldReplaceExistingContainers, shouldRestoreMappings *bool,
	timeoutSecs *uint32,
	containerAccountIds []valueObject.AccountId,
	containerIds []valueObject.ContainerId,
	exceptContainerAccountIds []valueObject.AccountId,
	exceptContainerIds []valueObject.ContainerId,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) RestoreBackupTask {
	return RestoreBackupTask{
		TaskId:                          taskId,
		ArchiveId:                       archiveId,
		ShouldReplaceExistingContainers: shouldReplaceExistingContainers,
		ShouldRestoreMappings:           shouldRestoreMappings,
		TimeoutSecs:                     timeoutSecs,
		ContainerAccountIds:             containerAccountIds,
		ContainerIds:                    containerIds,
		ExceptContainerAccountIds:       exceptContainerAccountIds,
		ExceptContainerIds:              exceptContainerIds,
		OperatorAccountId:               operatorAccountId,
		OperatorIpAddress:               operatorIpAddress,
	}
}
