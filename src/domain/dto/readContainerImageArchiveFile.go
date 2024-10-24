package dto

import (
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type ReadContainerImageArchiveFile struct {
	AccountId valueObject.AccountId        `json:"accountId"`
	ImageId   valueObject.ContainerImageId `json:"imageId"`
}

func NewReadContainerImageArchiveFile(
	accountId valueObject.AccountId,
	imageId valueObject.ContainerImageId,
) ReadContainerImageArchiveFile {
	return ReadContainerImageArchiveFile{
		AccountId: accountId,
		ImageId:   imageId,
	}
}
