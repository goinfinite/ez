package dto

import (
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type DeleteBackupDestination struct {
	DestinationId     valueObject.BackupDestinationId `json:"destinationId"`
	AccountId         valueObject.AccountId           `json:"accountId"`
	OperatorAccountId valueObject.AccountId           `json:"-"`
	OperatorIpAddress valueObject.IpAddress           `json:"-"`
}

func NewDeleteBackupDestination(
	destinationId valueObject.BackupDestinationId,
	accountId valueObject.AccountId,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) DeleteBackupDestination {
	return DeleteBackupDestination{
		AccountId:         accountId,
		DestinationId:     destinationId,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
