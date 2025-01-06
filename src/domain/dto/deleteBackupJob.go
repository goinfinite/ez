package dto

import (
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type DeleteBackupJob struct {
	JobId             valueObject.BackupJobId `json:"jobId"`
	AccountId         valueObject.AccountId   `json:"accountId"`
	OperatorAccountId valueObject.AccountId   `json:"-"`
	OperatorIpAddress valueObject.IpAddress   `json:"-"`
}

func NewDeleteBackupJob(
	jobId valueObject.BackupJobId,
	accountId valueObject.AccountId,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) DeleteBackupJob {
	return DeleteBackupJob{
		AccountId:         accountId,
		JobId:             jobId,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
