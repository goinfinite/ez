package dbModel

import (
	"time"

	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type BackupJob struct {
	ID                        uint64 `gorm:"primarykey"`
	AccountID                 uint64 `gorm:"not null"`
	JobStatus                 bool   `gorm:"not null"`
	JobDescription            *string
	DestinationIds            []uint64 `gorm:"serializer:json"`
	BackupType                string   `gorm:"not null"`
	BackupSchedule            string   `gorm:"not null"`
	ArchiveCompressionFormat  string   `gorm:"not null"`
	TimeoutSecs               uint64   `gorm:"not null"`
	MaxTaskRetentionCount     *uint16
	MaxTaskRetentionDays      *uint16
	MaxConcurrentCpuCores     *uint16
	ContainerAccountIds       []uint64 `gorm:"serializer:json"`
	ContainerIds              []uint64 `gorm:"serializer:json"`
	IgnoreContainerAccountIds []uint64 `gorm:"serializer:json"`
	IgnoreContainerIds        []uint64 `gorm:"serializer:json"`
	TasksCount                *uint16
	TotalSpaceUsageBytes      *uint64
	LastRunAt                 *time.Time
	LastRunStatus             *string
	NextRunAt                 *time.Time
	CreatedAt                 time.Time `gorm:"not null"`
	UpdatedAt                 time.Time `gorm:"not null"`
}

func (model BackupJob) TableName() string {
	return "backup_jobs"
}

func (model BackupJob) ToEntity() (jobEntity entity.BackupJob, err error) {
	jobId, err := valueObject.NewBackupJobId(model.ID)
	if err != nil {
		return jobEntity, err
	}

	accountId, err := valueObject.NewAccountId(model.AccountID)
	if err != nil {
		return jobEntity, err
	}

	var jobDescriptionPtr *valueObject.BackupJobDescription
	if model.JobDescription != nil {
		jobDescription, err := valueObject.NewBackupJobDescription(*model.JobDescription)
		if err != nil {
			return jobEntity, err
		}
		jobDescriptionPtr = &jobDescription
	}

	destinationIds := []valueObject.BackupDestinationId{}
	for _, destinationId := range model.DestinationIds {
		destinationId, err := valueObject.NewBackupDestinationId(destinationId)
		if err != nil {
			return jobEntity, err
		}
		destinationIds = append(destinationIds, destinationId)
	}

	backupType, err := valueObject.NewBackupJobType(model.BackupType)
	if err != nil {
		return jobEntity, err
	}

	backupSchedule, err := valueObject.NewCronSchedule(model.BackupSchedule)
	if err != nil {
		return jobEntity, err
	}

	archiveCompressionFormat, err := valueObject.NewCompressionFormat(model.ArchiveCompressionFormat)
	if err != nil {
		return jobEntity, err
	}

	containerAccountIds := []valueObject.AccountId{}
	for _, containerAccountId := range model.ContainerAccountIds {
		containerAccountId, err := valueObject.NewAccountId(containerAccountId)
		if err != nil {
			return jobEntity, err
		}
		containerAccountIds = append(containerAccountIds, containerAccountId)
	}

	containerIds := []valueObject.ContainerId{}
	for _, containerId := range model.ContainerIds {
		containerId, err := valueObject.NewContainerId(containerId)
		if err != nil {
			return jobEntity, err
		}
		containerIds = append(containerIds, containerId)
	}

	ignoreContainerAccountIds := []valueObject.AccountId{}
	for _, ignoreContainerAccountId := range model.IgnoreContainerAccountIds {
		ignoreContainerAccountId, err := valueObject.NewAccountId(ignoreContainerAccountId)
		if err != nil {
			return jobEntity, err
		}
		ignoreContainerAccountIds = append(ignoreContainerAccountIds, ignoreContainerAccountId)
	}

	ignoreContainerIds := []valueObject.ContainerId{}
	for _, ignoreContainerId := range model.IgnoreContainerIds {
		ignoreContainerId, err := valueObject.NewContainerId(ignoreContainerId)
		if err != nil {
			return jobEntity, err
		}
		ignoreContainerIds = append(ignoreContainerIds, ignoreContainerId)
	}

	var totalSpaceUsageBytesPtr *valueObject.Byte
	if model.TotalSpaceUsageBytes != nil {
		totalSpaceUsageBytes, err := valueObject.NewByte(*model.TotalSpaceUsageBytes)
		if err != nil {
			return jobEntity, err
		}
		totalSpaceUsageBytesPtr = &totalSpaceUsageBytes
	}

	var lastRunAtPtr *valueObject.UnixTime
	if model.LastRunAt != nil {
		lastRunAt := valueObject.NewUnixTimeWithGoTime(*model.LastRunAt)
		lastRunAtPtr = &lastRunAt
	}

	var lastRunStatusPtr *valueObject.BackupTaskStatus
	if model.LastRunStatus != nil {
		lastRunStatus, err := valueObject.NewBackupTaskStatus(*model.LastRunStatus)
		if err != nil {
			return jobEntity, err
		}
		lastRunStatusPtr = &lastRunStatus
	}

	var nextRunAtPtr *valueObject.UnixTime
	if model.NextRunAt != nil {
		nextRunAt := valueObject.NewUnixTimeWithGoTime(*model.NextRunAt)
		nextRunAtPtr = &nextRunAt
	}

	return entity.NewBackupJob(
		jobId, accountId, model.JobStatus, jobDescriptionPtr, destinationIds,
		backupType, backupSchedule, &archiveCompressionFormat, &model.TimeoutSecs,
		model.MaxTaskRetentionCount, model.MaxTaskRetentionDays,
		model.MaxConcurrentCpuCores, containerAccountIds, containerIds,
		ignoreContainerAccountIds, ignoreContainerIds, model.TasksCount,
		totalSpaceUsageBytesPtr, lastRunAtPtr, lastRunStatusPtr,
		nextRunAtPtr, valueObject.NewUnixTimeWithGoTime(model.CreatedAt),
		valueObject.NewUnixTimeWithGoTime(model.UpdatedAt),
	)
}
