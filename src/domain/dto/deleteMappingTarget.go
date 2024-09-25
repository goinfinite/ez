package dto

import (
	"github.com/speedianet/control/src/domain/valueObject"
)

type DeleteMappingTarget struct {
	AccountId         valueObject.AccountId       `json:"accountId"`
	MappingId         valueObject.MappingId       `json:"mappingId"`
	TargetId          valueObject.MappingTargetId `json:"targetId"`
	OperatorAccountId valueObject.AccountId       `json:"-"`
	OperatorIpAddress valueObject.IpAddress       `json:"-"`
}

func NewDeleteMappingTarget(
	accountId valueObject.AccountId,
	mappingId valueObject.MappingId,
	targetId valueObject.MappingTargetId,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) DeleteMappingTarget {
	return DeleteMappingTarget{
		AccountId:         accountId,
		MappingId:         mappingId,
		TargetId:          targetId,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
