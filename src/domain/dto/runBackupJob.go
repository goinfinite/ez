package dto

import (
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type RunBackupJob struct {
	JobId             valueObject.BackupJobId `json:"jobId"`
	AccountId         valueObject.AccountId   `json:"accountId"`
	OperatorAccountId valueObject.AccountId   `json:"-"`
	OperatorIpAddress valueObject.IpAddress   `json:"-"`
}

func NewRunBackupJob(
	jobId valueObject.BackupJobId,
	accountId, operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) RunBackupJob {
	return RunBackupJob{
		JobId:             jobId,
		AccountId:         accountId,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
