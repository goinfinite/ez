package entity

import (
	"errors"

	"github.com/goinfinite/ez/src/domain/valueObject"
)

type BackupJob struct {
	JobId                     valueObject.BackupJobId             `json:"jobId"`
	AccountId                 valueObject.AccountId               `json:"accountId"`
	AccountUsername           valueObject.UnixUsername            `json:"accountUsername"`
	JobStatus                 bool                                `json:"jobStatus"`
	JobDescription            *valueObject.BackupJobDescription   `json:"jobDescription"`
	DestinationIds            []valueObject.BackupDestinationId   `json:"destinationIds"`
	RetentionStrategy         valueObject.BackupRetentionStrategy `json:"retentionStrategy"`
	BackupSchedule            valueObject.CronSchedule            `json:"backupSchedule"`
	ArchiveCompressionFormat  valueObject.CompressionFormat       `json:"archiveCompressionFormat"`
	TimeoutSecs               valueObject.TimeDuration            `json:"timeoutSecs"`
	MaxTaskRetentionCount     *uint16                             `json:"maxTaskRetentionCount"`
	MaxTaskRetentionDays      *uint16                             `json:"maxTaskRetentionDays"`
	MaxConcurrentCpuCores     *uint16                             `json:"maxConcurrentCpuCores"`
	ContainerAccountIds       []valueObject.AccountId             `json:"containerAccountIds"`
	ContainerIds              []valueObject.ContainerId           `json:"containerIds"`
	ExceptContainerAccountIds []valueObject.AccountId             `json:"exceptContainerAccountIds"`
	ExceptContainerIds        []valueObject.ContainerId           `json:"exceptContainerIds"`
	TasksCount                uint16                              `json:"tasksCount"`
	TotalSpaceUsageBytes      valueObject.Byte                    `json:"totalSpaceUsageBytes"`
	LastRunAt                 *valueObject.UnixTime               `json:"lastRunAt"`
	LastRunStatus             *valueObject.BackupTaskStatus       `json:"lastRunStatus"`
	NextRunAt                 *valueObject.UnixTime               `json:"nextRunAt"`
	CreatedAt                 valueObject.UnixTime                `json:"createdAt"`
	UpdatedAt                 valueObject.UnixTime                `json:"updatedAt"`
}

func NewBackupJob(
	jobId valueObject.BackupJobId,
	accountId valueObject.AccountId,
	accountUsername valueObject.UnixUsername,
	jobStatus bool,
	jobDescription *valueObject.BackupJobDescription,
	destinationIds []valueObject.BackupDestinationId,
	retentionStrategy valueObject.BackupRetentionStrategy,
	backupSchedule valueObject.CronSchedule,
	archiveCompressionFormatPtr *valueObject.CompressionFormat,
	timeoutSecsPtr *valueObject.TimeDuration,
	maxTaskRetentionCount, maxTaskRetentionDays, maxConcurrentCpuCores *uint16,
	containerAccountIds []valueObject.AccountId,
	containerIds []valueObject.ContainerId,
	exceptContainerAccountIds []valueObject.AccountId,
	exceptContainerIds []valueObject.ContainerId,
	tasksCount uint16,
	totalSpaceUsageBytes valueObject.Byte,
	lastRunAt *valueObject.UnixTime,
	lastRunStatus *valueObject.BackupTaskStatus,
	nextRunAt *valueObject.UnixTime,
	createdAt, updatedAt valueObject.UnixTime,
) (backupJob BackupJob, err error) {
	if len(destinationIds) == 0 {
		return backupJob, errors.New("BackupJobMustHaveAtLeastOneDestination")
	}

	if retentionStrategy == valueObject.BackupRetentionStrategyIncremental {
		return backupJob, errors.New("IncrementalBackupJobNotSupportedYet")
	}

	archiveCompressionFormat := valueObject.CompressionFormatBrotli
	if archiveCompressionFormatPtr != nil {
		archiveCompressionFormat = *archiveCompressionFormatPtr
	}

	timeoutSecs := valueObject.TimeDuration(uint64(8 * 60 * 60))
	if timeoutSecsPtr != nil {
		timeoutSecs = *timeoutSecsPtr
	}

	return BackupJob{
		JobId:                     jobId,
		AccountId:                 accountId,
		AccountUsername:           accountUsername,
		JobStatus:                 jobStatus,
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
		TasksCount:                tasksCount,
		TotalSpaceUsageBytes:      totalSpaceUsageBytes,
		LastRunAt:                 lastRunAt,
		LastRunStatus:             lastRunStatus,
		NextRunAt:                 nextRunAt,
		CreatedAt:                 createdAt,
		UpdatedAt:                 updatedAt,
	}, nil
}
