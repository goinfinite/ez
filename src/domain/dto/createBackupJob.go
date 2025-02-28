package dto

import (
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type CreateBackupJob struct {
	AccountId                 valueObject.AccountId                `json:"accountId"`
	JobDescription            *valueObject.BackupJobDescription    `json:"jobDescription"`
	DestinationIds            []valueObject.BackupDestinationId    `json:"destinationIds"`
	RetentionStrategy         *valueObject.BackupRetentionStrategy `json:"retentionStrategy"`
	BackupSchedule            valueObject.CronSchedule             `json:"backupSchedule"`
	ArchiveCompressionFormat  *valueObject.CompressionFormat       `json:"archiveCompressionFormat,omitempty"`
	TimeoutSecs               *valueObject.TimeDuration            `json:"timeoutSecs,omitempty"`
	MaxTaskRetentionCount     *uint16                              `json:"maxTaskRetentionCount,omitempty"`
	MaxTaskRetentionDays      *uint16                              `json:"maxTaskRetentionDays,omitempty"`
	MaxConcurrentCpuCores     *uint16                              `json:"maxConcurrentCpuCores,omitempty"`
	ContainerAccountIds       []valueObject.AccountId              `json:"containerAccountIds"`
	ContainerIds              []valueObject.ContainerId            `json:"containerIds"`
	ExceptContainerAccountIds []valueObject.AccountId              `json:"exceptContainerAccountIds"`
	ExceptContainerIds        []valueObject.ContainerId            `json:"exceptContainerIds"`
	OperatorAccountId         valueObject.AccountId                `json:"-"`
	OperatorIpAddress         valueObject.IpAddress                `json:"-"`
}

func NewCreateBackupJob(
	accountId valueObject.AccountId,
	jobDescription *valueObject.BackupJobDescription,
	destinationIds []valueObject.BackupDestinationId,
	retentionStrategy *valueObject.BackupRetentionStrategy,
	backupSchedule valueObject.CronSchedule,
	archiveCompressionFormat *valueObject.CompressionFormat,
	timeoutSecs *valueObject.TimeDuration,
	maxTaskRetentionCount *uint16,
	maxTaskRetentionDays *uint16,
	maxConcurrentCpuCores *uint16,
	containerAccountIds []valueObject.AccountId,
	containerIds []valueObject.ContainerId,
	exceptContainerAccountIds []valueObject.AccountId,
	exceptContainerIds []valueObject.ContainerId,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) CreateBackupJob {
	return CreateBackupJob{
		AccountId:                 accountId,
		JobDescription:            jobDescription,
		DestinationIds:            destinationIds,
		RetentionStrategy:         retentionStrategy,
		BackupSchedule:            backupSchedule,
		ArchiveCompressionFormat:  archiveCompressionFormat,
		TimeoutSecs:               timeoutSecs,
		MaxTaskRetentionCount:     maxTaskRetentionCount,
		MaxTaskRetentionDays:      maxTaskRetentionDays,
		MaxConcurrentCpuCores:     maxConcurrentCpuCores,
		ContainerAccountIds:       containerAccountIds,
		ContainerIds:              containerIds,
		ExceptContainerAccountIds: exceptContainerAccountIds,
		ExceptContainerIds:        exceptContainerIds,
		OperatorAccountId:         operatorAccountId,
		OperatorIpAddress:         operatorIpAddress,
	}
}
