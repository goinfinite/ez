package dto

import "github.com/speedianet/control/src/domain/valueObject"

type CreateMappingTarget struct {
	AccountId         valueObject.AccountId   `json:"accountId"`
	MappingId         valueObject.MappingId   `json:"mappingId"`
	ContainerId       valueObject.ContainerId `json:"containerId"`
	OperatorAccountId valueObject.AccountId   `json:"-"`
	OperatorIpAddress valueObject.IpAddress   `json:"-"`
}

func NewCreateMappingTarget(
	accountId valueObject.AccountId,
	mappingId valueObject.MappingId,
	containerId valueObject.ContainerId,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) CreateMappingTarget {
	return CreateMappingTarget{
		AccountId:         accountId,
		MappingId:         mappingId,
		ContainerId:       containerId,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
