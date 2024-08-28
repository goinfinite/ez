package dto

import (
	"github.com/speedianet/control/src/domain/valueObject"
)

type DeleteContainerImage struct {
	AccountId         valueObject.AccountId        `json:"accountId"`
	ImageId           valueObject.ContainerImageId `json:"imageId"`
	OperatorAccountId valueObject.AccountId        `json:"-"`
	IpAddress         valueObject.IpAddress        `json:"-"`
}

func NewDeleteContainerImage(
	accountId valueObject.AccountId,
	imageId valueObject.ContainerImageId,
	operatorAccountId valueObject.AccountId,
	ipAddress valueObject.IpAddress,
) DeleteContainerImage {
	return DeleteContainerImage{
		AccountId:         accountId,
		ImageId:           imageId,
		OperatorAccountId: operatorAccountId,
		IpAddress:         ipAddress,
	}
}
