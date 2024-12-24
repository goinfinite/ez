package entity

import (
	"errors"

	"github.com/goinfinite/ez/src/domain/valueObject"
)

type BackupJob struct {
	JobId                     valueObject.BackupJobId           `json:"jobId"`
	AccountId                 valueObject.AccountId             `json:"accountId"`
	JobStatus                 bool                              `json:"jobsStatus"`
	JobDescription            *valueObject.BackupJobDescription `json:"jobDescription"`
	DestinationIds            []valueObject.BackupDestinationId `json:"destinationIds"`
	BackupType                valueObject.BackupJobType         `json:"backupType"`
	BackupSchedule            valueObject.CronSchedule          `json:"backupSchedule"`
	ArchiveCompressionFormat  *valueObject.CompressionFormat    `json:"archiveCompressionFormat"`
	TimeoutSecs               *uint64                           `json:"timeoutSecs"`
	MaxTaskRetentionCount     *uint16                           `json:"maxTaskRetentionCount"`
	MaxTaskRetentionDays      *uint16                           `json:"maxTaskRetentionDays"`
	MaxConcurrentCpuCores     *uint16                           `json:"maxConcurrentCpuCores"`
	ContainerAccountIds       []valueObject.AccountId           `json:"containerAccountIds"`
	ContainerIds              []valueObject.ContainerId         `json:"containerIds"`
	IgnoreContainerAccountIds []valueObject.AccountId           `json:"ignoreContainerAccountIds"`
	IgnoreContainerIds        []valueObject.ContainerId         `json:"ignoreContainerIds"`
	TasksCount                *uint16                           `json:"tasksCount"`
	TotalSpaceUsageBytes      *valueObject.Byte                 `json:"totalSpaceUsageBytes"`
	LastRunAt                 *valueObject.UnixTime             `json:"lastRunAt"`
	LastRunStatus             *valueObject.BackupTaskStatus     `json:"lastRunStatus"`
	NextRunAt                 *valueObject.UnixTime             `json:"nextRunAt"`
	CreatedAt                 valueObject.UnixTime              `json:"createdAt"`
	UpdatedAt                 valueObject.UnixTime              `json:"updatedAt"`
}

func NewBackupJob(
	jobId valueObject.BackupJobId,
	accountId valueObject.AccountId,
	jobStatus bool,
	jobDescription *valueObject.BackupJobDescription,
	destinationIds []valueObject.BackupDestinationId,
	backupType valueObject.BackupJobType,
	backupSchedule valueObject.CronSchedule,
	archiveCompressionFormatPtr *valueObject.CompressionFormat,
	timeoutSecsPtr *uint64,
	maxTaskRetentionCount, maxTaskRetentionDays, maxConcurrentCpuCores *uint16,
	containerAccountIds []valueObject.AccountId,
	containerIds []valueObject.ContainerId,
	ignoreContainerAccountIds []valueObject.AccountId,
	ignoreContainerIds []valueObject.ContainerId,
	tasksCount *uint16,
	totalSpaceUsageBytes *valueObject.Byte,
	lastRunAt *valueObject.UnixTime,
	lastRunStatus *valueObject.BackupTaskStatus,
	nextRunAt *valueObject.UnixTime,
	createdAt, updatedAt valueObject.UnixTime,
) (backupJob BackupJob, err error) {
	if len(destinationIds) == 0 {
		return backupJob, errors.New("BackupJobMustHaveAtLeastOneDestination")
	}

	if backupType == valueObject.BackupJobTypeIncremental {
		return backupJob, errors.New("IncrementalBackupJobNotSupportedYet")
	}

	if archiveCompressionFormatPtr == nil {
		archiveCompressionFormat := valueObject.CompressionFormatBrotli
		archiveCompressionFormatPtr = &archiveCompressionFormat
	}

	if timeoutSecsPtr == nil {
		timeoutSecs := uint64(8 * 60 * 60)
		timeoutSecsPtr = &timeoutSecs
	}

	return BackupJob{
		JobId:                     jobId,
		AccountId:                 accountId,
		JobStatus:                 jobStatus,
		JobDescription:            jobDescription,
		DestinationIds:            destinationIds,
		BackupType:                backupType,
		BackupSchedule:            backupSchedule,
		ArchiveCompressionFormat:  archiveCompressionFormatPtr,
		TimeoutSecs:               timeoutSecsPtr,
		MaxTaskRetentionCount:     maxTaskRetentionCount,
		MaxTaskRetentionDays:      maxTaskRetentionDays,
		MaxConcurrentCpuCores:     maxConcurrentCpuCores,
		ContainerAccountIds:       containerAccountIds,
		ContainerIds:              containerIds,
		IgnoreContainerAccountIds: ignoreContainerAccountIds,
		IgnoreContainerIds:        ignoreContainerIds,
		TasksCount:                tasksCount,
		TotalSpaceUsageBytes:      totalSpaceUsageBytes,
		LastRunAt:                 lastRunAt,
		LastRunStatus:             lastRunStatus,
		NextRunAt:                 nextRunAt,
		CreatedAt:                 createdAt,
		UpdatedAt:                 updatedAt,
	}, nil
}
