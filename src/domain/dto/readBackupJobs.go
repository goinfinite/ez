package dto

import (
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type ReadBackupJobsRequest struct {
	Pagination               Pagination                           `json:"pagination"`
	JobId                    *valueObject.BackupJobId             `json:"jobId"`
	JobStatus                *bool                                `json:"jobStatus"`
	AccountId                *valueObject.AccountId               `json:"accountId"`
	DestinationId            *valueObject.BackupDestinationId     `json:"destinationId"`
	RetentionStrategy        *valueObject.BackupRetentionStrategy `json:"retentionStrategy"`
	ArchiveCompressionFormat *valueObject.CompressionFormat       `json:"archiveCompressionFormat"`
	LastRunStatus            *valueObject.BackupTaskStatus        `json:"lastRunStatus"`
	LastRunBeforeAt          *valueObject.UnixTime                `json:"lastRunBeforeAt"`
	LastRunAfterAt           *valueObject.UnixTime                `json:"lastRunAfterAt"`
	NextRunBeforeAt          *valueObject.UnixTime                `json:"nextRunBeforeAt"`
	NextRunAfterAt           *valueObject.UnixTime                `json:"nextRunAfterAt"`
	CreatedBeforeAt          *valueObject.UnixTime                `json:"createdBeforeAt"`
	CreatedAfterAt           *valueObject.UnixTime                `json:"createdAfterAt"`
}

type ReadBackupJobsResponse struct {
	Pagination Pagination         `json:"pagination"`
	Jobs       []entity.BackupJob `json:"jobs"`
}
