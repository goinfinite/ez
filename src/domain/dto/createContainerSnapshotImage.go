package dto

import (
	"github.com/speedianet/control/src/domain/valueObject"
)

type CreateContainerSnapshotImage struct {
	AccountId           valueObject.AccountId          `json:"accountId"`
	ContainerId         valueObject.ContainerId        `json:"containerId"`
	ShouldCreateArchive *bool                          `json:"shouldCreateArchive"`
	CompressionFormat   *valueObject.CompressionFormat `json:"compressionFormat"`
	OperatorAccountId   valueObject.AccountId          `json:"-"`
	OperatorIpAddress   valueObject.IpAddress          `json:"-"`
}

func NewCreateContainerSnapshotImage(
	accountId valueObject.AccountId,
	containerId valueObject.ContainerId,
	shouldCreateArchive *bool,
	compressionFormat *valueObject.CompressionFormat,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) CreateContainerSnapshotImage {
	return CreateContainerSnapshotImage{
		AccountId:           accountId,
		ContainerId:         containerId,
		ShouldCreateArchive: shouldCreateArchive,
		CompressionFormat:   compressionFormat,
		OperatorAccountId:   operatorAccountId,
		OperatorIpAddress:   operatorIpAddress,
	}
}
