package repository

import (
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
)

type ContainerQueryRepo interface {
	Read(dto.ReadContainersRequest) (dto.ReadContainersResponse, error)
	ReadFirst(dto.ReadContainersRequest) (entity.Container, error)
	ReadFirstWithMetrics(dto.ReadContainersRequest) (dto.ContainerWithMetrics, error)
}
