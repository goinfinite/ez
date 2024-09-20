package dto

import (
	"github.com/speedianet/control/src/domain/valueObject"
)

type CreateContainerSnapshotImage struct {
	AccountId                valueObject.AccountId          `json:"accountId"`
	ContainerId              valueObject.ContainerId        `json:"containerId"`
	ShouldCreateArchive      *bool                          `json:"shouldCreateArchive"`
	ArchiveCompressionFormat *valueObject.CompressionFormat `json:"archiveCompressionFormat"`
	ShouldDiscardImage       *bool                          `json:"shouldDiscardImage"`
	OperatorAccountId        valueObject.AccountId          `json:"-"`
	OperatorIpAddress        valueObject.IpAddress          `json:"-"`
}

func NewCreateContainerSnapshotImage(
	accountId valueObject.AccountId,
	containerId valueObject.ContainerId,
	shouldCreateArchive *bool,
	archiveCompressionFormat *valueObject.CompressionFormat,
	shouldDiscardImage *bool,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) CreateContainerSnapshotImage {
	return CreateContainerSnapshotImage{
		AccountId:                accountId,
		ContainerId:              containerId,
		ShouldCreateArchive:      shouldCreateArchive,
		ArchiveCompressionFormat: archiveCompressionFormat,
		ShouldDiscardImage:       shouldDiscardImage,
		OperatorAccountId:        operatorAccountId,
		OperatorIpAddress:        operatorIpAddress,
	}
}
