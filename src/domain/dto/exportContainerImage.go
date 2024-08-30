package dto

import (
	"github.com/speedianet/control/src/domain/valueObject"
)

type ExportContainerImage struct {
	AccountId         valueObject.AccountId        `json:"accountId"`
	ImageId           valueObject.ContainerImageId `json:"imageId"`
	OperatorAccountId valueObject.AccountId        `json:"-"`
	IpAddress         valueObject.IpAddress        `json:"-"`
}

func NewExportContainerImage(
	accountId valueObject.AccountId,
	imageId valueObject.ContainerImageId,
	operatorAccountId valueObject.AccountId,
	ipAddress valueObject.IpAddress,
) ExportContainerImage {
	return ExportContainerImage{
		AccountId:         accountId,
		ImageId:           imageId,
		OperatorAccountId: operatorAccountId,
		IpAddress:         ipAddress,
	}
}
