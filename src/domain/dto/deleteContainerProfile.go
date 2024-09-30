package dto

import (
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type DeleteContainerProfile struct {
	AccountId         valueObject.AccountId          `json:"accountId"`
	ProfileId         valueObject.ContainerProfileId `json:"profileId"`
	OperatorAccountId valueObject.AccountId          `json:"-"`
	OperatorIpAddress valueObject.IpAddress          `json:"-"`
}

func NewDeleteContainerProfile(
	accountId valueObject.AccountId,
	profileId valueObject.ContainerProfileId,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) DeleteContainerProfile {
	return DeleteContainerProfile{
		AccountId:         accountId,
		ProfileId:         profileId,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
