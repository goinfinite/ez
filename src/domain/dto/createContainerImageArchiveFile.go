package dto

import (
	"github.com/speedianet/control/src/domain/valueObject"
)

type CreateContainerImageArchiveFile struct {
	AccountId         valueObject.AccountId        `json:"accountId"`
	ImageId           valueObject.ContainerImageId `json:"imageId"`
	OperatorAccountId valueObject.AccountId        `json:"-"`
	OperatorIpAddress valueObject.IpAddress        `json:"-"`
}

func NewCreateContainerImageArchiveFile(
	accountId valueObject.AccountId,
	imageId valueObject.ContainerImageId,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) CreateContainerImageArchiveFile {
	return CreateContainerImageArchiveFile{
		AccountId:         accountId,
		ImageId:           imageId,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
