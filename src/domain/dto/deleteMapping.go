package dto

import (
	"github.com/speedianet/control/src/domain/valueObject"
)

type DeleteMapping struct {
	AccountId         valueObject.AccountId `json:"accountId"`
	MappingId         valueObject.MappingId `json:"mappingId"`
	OperatorAccountId valueObject.AccountId `json:"-"`
	OperatorIpAddress valueObject.IpAddress `json:"-"`
}

func NewDeleteMapping(
	accountId valueObject.AccountId,
	mappingId valueObject.MappingId,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) DeleteMapping {
	return DeleteMapping{
		AccountId:         accountId,
		MappingId:         mappingId,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
