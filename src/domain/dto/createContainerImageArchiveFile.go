package dto

import (
	"github.com/speedianet/control/src/domain/valueObject"
)

type CreateContainerImageArchiveFile struct {
	AccountId         valueObject.AccountId          `json:"accountId"`
	ImageId           valueObject.ContainerImageId   `json:"imageId"`
	CompressionFormat *valueObject.CompressionFormat `json:"compressionFormat"`
	OperatorAccountId valueObject.AccountId          `json:"-"`
	OperatorIpAddress valueObject.IpAddress          `json:"-"`
}

func NewCreateContainerImageArchiveFile(
	accountId valueObject.AccountId,
	imageId valueObject.ContainerImageId,
	compressionFormat *valueObject.CompressionFormat,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) CreateContainerImageArchiveFile {
	return CreateContainerImageArchiveFile{
		AccountId:         accountId,
		ImageId:           imageId,
		CompressionFormat: compressionFormat,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
