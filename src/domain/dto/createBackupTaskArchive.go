package dto

import (
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type CreateBackupTaskArchive struct {
	TaskId                    valueObject.BackupTaskId  `json:"taskId"`
	TimeoutSecs               *uint32                   `json:"timeoutSecs,omitempty"`
	ContainerAccountIds       []valueObject.AccountId   `json:"containerAccountIds"`
	ContainerIds              []valueObject.ContainerId `json:"containerIds"`
	ExceptContainerAccountIds []valueObject.AccountId   `json:"exceptContainerAccountIds"`
	ExceptContainerIds        []valueObject.ContainerId `json:"exceptContainerIds"`
	OperatorAccountId         valueObject.AccountId     `json:"-"`
	OperatorIpAddress         valueObject.IpAddress     `json:"-"`
}

func NewCreateBackupTaskArchive(
	taskId valueObject.BackupTaskId,
	timeoutSecs *uint32,
	containerAccountIds []valueObject.AccountId,
	containerIds []valueObject.ContainerId,
	exceptContainerAccountIds []valueObject.AccountId,
	exceptContainerIds []valueObject.ContainerId,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) CreateBackupTaskArchive {
	return CreateBackupTaskArchive{
		TaskId:                    taskId,
		TimeoutSecs:               timeoutSecs,
		ContainerAccountIds:       containerAccountIds,
		ContainerIds:              containerIds,
		ExceptContainerAccountIds: exceptContainerAccountIds,
		ExceptContainerIds:        exceptContainerIds,
		OperatorAccountId:         operatorAccountId,
		OperatorIpAddress:         operatorIpAddress,
	}
}
