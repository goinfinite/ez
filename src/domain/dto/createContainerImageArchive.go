package dto

import (
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type CreateContainerImageArchive struct {
	AccountId         valueObject.AccountId          `json:"accountId"`
	ImageId           valueObject.ContainerImageId   `json:"imageId"`
	CompressionFormat *valueObject.CompressionFormat `json:"compressionFormat"`
	DestinationPath   *valueObject.UnixFilePath      `json:"-"`
	OperatorAccountId valueObject.AccountId          `json:"-"`
	OperatorIpAddress valueObject.IpAddress          `json:"-"`
}

func NewCreateContainerImageArchive(
	accountId valueObject.AccountId,
	imageId valueObject.ContainerImageId,
	compressionFormat *valueObject.CompressionFormat,
	destinationPath *valueObject.UnixFilePath,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) CreateContainerImageArchive {
	return CreateContainerImageArchive{
		AccountId:         accountId,
		ImageId:           imageId,
		CompressionFormat: compressionFormat,
		DestinationPath:   destinationPath,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
