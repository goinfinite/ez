package dto

import (
	"github.com/speedianet/control/src/domain/valueObject"
)

type DeleteContainerImageArchiveFile struct {
	AccountId         valueObject.AccountId        `json:"accountId"`
	ImageId           valueObject.ContainerImageId `json:"imageId"`
	OperatorAccountId valueObject.AccountId        `json:"-"`
	IpAddress         valueObject.IpAddress        `json:"-"`
}

func NewDeleteContainerImageArchiveFile(
	accountId valueObject.AccountId,
	imageId valueObject.ContainerImageId,
	operatorAccountId valueObject.AccountId,
	ipAddress valueObject.IpAddress,
) DeleteContainerImageArchiveFile {
	return DeleteContainerImageArchiveFile{
		AccountId:         accountId,
		ImageId:           imageId,
		OperatorAccountId: operatorAccountId,
		IpAddress:         ipAddress,
	}
}
