package dto

import (
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type UpdateBackupJob struct {
	JobId                     valueObject.BackupJobId           `json:"jobId"`
	AccountId                 valueObject.AccountId             `json:"accountId"`
	JobStatus                 *bool                             `json:"jobStatus,omitempty"`
	JobDescription            *valueObject.BackupJobDescription `json:"jobDescription,omitempty"`
	DestinationIds            []valueObject.BackupDestinationId `json:"destinationIds,omitempty"`
	BackupSchedule            *valueObject.CronSchedule         `json:"backupSchedule,omitempty"`
	TimeoutSecs               *valueObject.TimeDuration         `json:"timeoutSecs,omitempty"`
	MaxTaskRetentionCount     *uint16                           `json:"maxTaskRetentionCount,omitempty"`
	MaxTaskRetentionDays      *uint16                           `json:"maxTaskRetentionDays,omitempty"`
	MaxConcurrentCpuCores     *uint16                           `json:"maxConcurrentCpuCores,omitempty"`
	ContainerAccountIds       []valueObject.AccountId           `json:"containerAccountIds,omitempty"`
	ContainerIds              []valueObject.ContainerId         `json:"containerIds,omitempty"`
	ExceptContainerAccountIds []valueObject.AccountId           `json:"exceptContainerAccountIds,omitempty"`
	ExceptContainerIds        []valueObject.ContainerId         `json:"exceptContainerIds,omitempty"`
	TasksCount                *uint16                           `json:"-"`
	TotalSpaceUsageBytes      *valueObject.Byte                 `json:"-"`
	LastRunAt                 *valueObject.UnixTime             `json:"-"`
	LastRunStatus             *valueObject.BackupTaskStatus     `json:"-"`
	NextRunAt                 *valueObject.UnixTime             `json:"-"`
	OperatorAccountId         valueObject.AccountId             `json:"-"`
	OperatorIpAddress         valueObject.IpAddress             `json:"-"`
}
