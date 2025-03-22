package dto

import (
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type ReadContainerImageArchivesRequest struct {
	Pagination        Pagination                    `json:"pagination"`
	ImageId           *valueObject.ContainerImageId `json:"imageId,omitempty"`
	AccountId         *valueObject.AccountId        `json:"accountId,omitempty"`
	CreatedBeforeAt   *valueObject.UnixTime         `json:"createdBeforeAt,omitempty"`
	CreatedAfterAt    *valueObject.UnixTime         `json:"createdAfterAt,omitempty"`
	ArchivesDirectory *valueObject.UnixFilePath     `json:"-"`
}

type ReadContainerImageArchivesResponse struct {
	Pagination Pagination                     `json:"pagination"`
	Archives   []entity.ContainerImageArchive `json:"archives"`
}
