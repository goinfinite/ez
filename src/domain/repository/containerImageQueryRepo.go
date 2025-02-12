package repository

import (
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type ContainerImageQueryRepo interface {
	Read() ([]entity.ContainerImage, error)
	ReadById(
		accountId valueObject.AccountId,
		imageId valueObject.ContainerImageId,
	) (entity.ContainerImage, error)
	ReadArchives(
		dto.ReadContainerImageArchivesRequest,
	) (dto.ReadContainerImageArchivesResponse, error)
	ReadArchive(
		dto.ReadContainerImageArchive,
	) (entity.ContainerImageArchive, error)
}
