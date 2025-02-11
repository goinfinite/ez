package dto

import (
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type DeleteBackupTaskArchive struct {
	ArchiveId         valueObject.BackupTaskArchiveId `json:"archiveId"`
	OperatorAccountId valueObject.AccountId           `json:"-"`
	OperatorIpAddress valueObject.IpAddress           `json:"-"`
}

func NewDeleteBackupTaskArchive(
	archiveId valueObject.BackupTaskArchiveId,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) DeleteBackupTaskArchive {
	return DeleteBackupTaskArchive{
		ArchiveId:         archiveId,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
