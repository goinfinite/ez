package dto

import (
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type CreateContainerSnapshotImage struct {
	ContainerId              valueObject.ContainerId        `json:"containerId"`
	ShouldCreateArchive      *bool                          `json:"shouldCreateArchive"`
	ArchiveCompressionFormat *valueObject.CompressionFormat `json:"archiveCompressionFormat"`
	ShouldDiscardImage       *bool                          `json:"shouldDiscardImage"`
	OperatorAccountId        valueObject.AccountId          `json:"-"`
	OperatorIpAddress        valueObject.IpAddress          `json:"-"`
}

func NewCreateContainerSnapshotImage(
	containerId valueObject.ContainerId,
	shouldCreateArchive *bool,
	archiveCompressionFormat *valueObject.CompressionFormat,
	shouldDiscardImage *bool,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) CreateContainerSnapshotImage {
	return CreateContainerSnapshotImage{
		ContainerId:              containerId,
		ShouldCreateArchive:      shouldCreateArchive,
		ArchiveCompressionFormat: archiveCompressionFormat,
		ShouldDiscardImage:       shouldDiscardImage,
		OperatorAccountId:        operatorAccountId,
		OperatorIpAddress:        operatorIpAddress,
	}
}
