package dto

import (
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type ReadContainerImageArchive struct {
	AccountId valueObject.AccountId        `json:"accountId"`
	ImageId   valueObject.ContainerImageId `json:"imageId"`
}

func NewReadContainerImageArchive(
	accountId valueObject.AccountId,
	imageId valueObject.ContainerImageId,
) ReadContainerImageArchive {
	return ReadContainerImageArchive{
		AccountId: accountId,
		ImageId:   imageId,
	}
}
