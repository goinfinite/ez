package dto

import (
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type DeleteContainerImage struct {
	AccountId         valueObject.AccountId        `json:"accountId"`
	ImageId           valueObject.ContainerImageId `json:"imageId"`
	OperatorAccountId valueObject.AccountId        `json:"-"`
	OperatorIpAddress valueObject.IpAddress        `json:"-"`
}

func NewDeleteContainerImage(
	accountId valueObject.AccountId,
	imageId valueObject.ContainerImageId,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) DeleteContainerImage {
	return DeleteContainerImage{
		AccountId:         accountId,
		ImageId:           imageId,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
