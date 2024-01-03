package dto

import "github.com/speedianet/control/src/domain/valueObject"

type AddMappingTarget struct {
	AccountId valueObject.AccountId     `json:"accountId"`
	MappingId valueObject.MappingId     `json:"mappingId"`
	Target    valueObject.MappingTarget `json:"target"`
}

func NewAddMappingTarget(
	accountId valueObject.AccountId,
	mappingId valueObject.MappingId,
	target valueObject.MappingTarget,
) AddMappingTarget {
	return AddMappingTarget{
		AccountId: accountId,
		MappingId: mappingId,
		Target:    target,
	}
}
