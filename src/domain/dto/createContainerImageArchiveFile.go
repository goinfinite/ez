package dto

import (
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type CreateContainerImageArchiveFile struct {
	AccountId         valueObject.AccountId          `json:"accountId"`
	ImageId           valueObject.ContainerImageId   `json:"imageId"`
	CompressionFormat *valueObject.CompressionFormat `json:"compressionFormat"`
	DestinationPath   *valueObject.UnixFilePath      `json:"-"`
	OperatorAccountId valueObject.AccountId          `json:"-"`
	OperatorIpAddress valueObject.IpAddress          `json:"-"`
}

func NewCreateContainerImageArchiveFile(
	accountId valueObject.AccountId,
	imageId valueObject.ContainerImageId,
	compressionFormat *valueObject.CompressionFormat,
	destinationPath *valueObject.UnixFilePath,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) CreateContainerImageArchiveFile {
	return CreateContainerImageArchiveFile{
		AccountId:         accountId,
		ImageId:           imageId,
		CompressionFormat: compressionFormat,
		DestinationPath:   destinationPath,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
