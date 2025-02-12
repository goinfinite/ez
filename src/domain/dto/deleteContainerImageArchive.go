package dto

import (
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type DeleteContainerImageArchive struct {
	AccountId         valueObject.AccountId        `json:"accountId"`
	ImageId           valueObject.ContainerImageId `json:"imageId"`
	OperatorAccountId valueObject.AccountId        `json:"-"`
	OperatorIpAddress valueObject.IpAddress        `json:"-"`
}

func NewDeleteContainerImageArchive(
	accountId valueObject.AccountId,
	imageId valueObject.ContainerImageId,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) DeleteContainerImageArchive {
	return DeleteContainerImageArchive{
		AccountId:         accountId,
		ImageId:           imageId,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
