package dto

import (
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type ReadBackupTaskArchivesRequest struct {
	Pagination      Pagination                       `json:"pagination"`
	ArchiveId       *valueObject.BackupTaskArchiveId `json:"archiveId,omitempty"`
	AccountId       *valueObject.AccountId           `json:"accountId,omitempty"`
	TaskId          *valueObject.BackupTaskId        `json:"taskId,omitempty"`
	CreatedBeforeAt *valueObject.UnixTime            `json:"createdBeforeAt,omitempty"`
	CreatedAfterAt  *valueObject.UnixTime            `json:"createdAfterAt,omitempty"`
}

type ReadBackupTaskArchivesResponse struct {
	Pagination Pagination                 `json:"pagination"`
	Archives   []entity.BackupTaskArchive `json:"archives"`
}
