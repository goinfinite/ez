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
	TimeoutSecs               *uint64                              `json:"timeoutSecs,omitempty"`
	MaxTaskRetentionCount     *uint16                              `json:"maxTaskRetentionCount,omitempty"`
	MaxTaskRetentionDays      *uint16                              `json:"maxTaskRetentionDays,omitempty"`
	MaxConcurrentCpuCores     *uint16                              `json:"maxConcurrentCpuCores,omitempty"`
	ContainerAccountIds       []valueObject.AccountId              `json:"containerAccountIds"`
	ContainerIds              []valueObject.ContainerId            `json:"containerIds"`
	IgnoreContainerAccountIds []valueObject.AccountId              `json:"ignoreContainerAccountIds"`
	IgnoreContainerIds        []valueObject.ContainerId            `json:"ignoreContainerIds"`
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
	timeoutSecs *uint64,
	maxTaskRetentionCount *uint16,
	maxTaskRetentionDays *uint16,
	maxConcurrentCpuCores *uint16,
	containerAccountIds []valueObject.AccountId,
	containerIds []valueObject.ContainerId,
	ignoreContainerAccountIds []valueObject.AccountId,
	ignoreContainerIds []valueObject.ContainerId,
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
		IgnoreContainerAccountIds: ignoreContainerAccountIds,
		IgnoreContainerIds:        ignoreContainerIds,
		OperatorAccountId:         operatorAccountId,
		OperatorIpAddress:         operatorIpAddress,
	}
}
