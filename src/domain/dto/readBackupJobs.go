package dto

import (
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type ReadBackupJobsRequest struct {
	Pagination               Pagination                       `json:"pagination"`
	JobId                    *valueObject.BackupJobId         `json:"jobId,omitempty"`
	JobStatus                *bool                            `json:"jobStatus,omitempty"`
	AccountId                *valueObject.AccountId           `json:"accountId,omitempty"`
	DestinationId            *valueObject.BackupDestinationId `json:"destinationId,omitempty"`
	BackupType               *valueObject.BackupJobType       `json:"backupType,omitempty"`
	ArchiveCompressionFormat *valueObject.CompressionFormat   `json:"archiveCompressionFormat,omitempty"`
	LastRunStatus            *valueObject.BackupTaskStatus    `json:"lastRunStatus,omitempty"`
	LastRunBeforeAt          *valueObject.UnixTime            `json:"lastRunBeforeAt,omitempty"`
	LastRunAfterAt           *valueObject.UnixTime            `json:"lastRunAfterAt,omitempty"`
	NextRunBeforeAt          *valueObject.UnixTime            `json:"nextRunBeforeAt,omitempty"`
	NextRunAfterAt           *valueObject.UnixTime            `json:"nextRunAfterAt,omitempty"`
	CreatedBeforeAt          *valueObject.UnixTime            `json:"createdBeforeAt,omitempty"`
	CreatedAfterAt           *valueObject.UnixTime            `json:"createdAfterAt,omitempty"`
}

type ReadBackupJobsResponse struct {
	Pagination Pagination         `json:"pagination"`
	Jobs       []entity.BackupJob `json:"jobs"`
}
