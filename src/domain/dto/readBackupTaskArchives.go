package dto

import (
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type ReadBackupTaskArchivesRequest struct {
	Pagination      Pagination                       `json:"pagination"`
	ArchiveId       *valueObject.BackupTaskArchiveId `json:"archiveId"`
	AccountId       *valueObject.AccountId           `json:"accountId"`
	TaskId          *valueObject.BackupTaskId        `json:"taskId"`
	CreatedBeforeAt *valueObject.UnixTime            `json:"createdBeforeAt"`
	CreatedAfterAt  *valueObject.UnixTime            `json:"createdAfterAt"`
}

type ReadBackupTaskArchivesResponse struct {
	Pagination Pagination                 `json:"pagination"`
	Archives   []entity.BackupTaskArchive `json:"archives"`
}
