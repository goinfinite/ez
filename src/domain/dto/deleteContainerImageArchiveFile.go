package dto

import (
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type DeleteContainerImageArchiveFile struct {
	AccountId         valueObject.AccountId        `json:"accountId"`
	ImageId           valueObject.ContainerImageId `json:"imageId"`
	OperatorAccountId valueObject.AccountId        `json:"-"`
	OperatorIpAddress valueObject.IpAddress        `json:"-"`
}

func NewDeleteContainerImageArchiveFile(
	accountId valueObject.AccountId,
	imageId valueObject.ContainerImageId,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) DeleteContainerImageArchiveFile {
	return DeleteContainerImageArchiveFile{
		AccountId:         accountId,
		ImageId:           imageId,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
