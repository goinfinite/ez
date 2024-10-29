package repository

import (
	"github.com/goinfinite/ez/src/domain/dto"
)

type ContainerQueryRepo interface {
	Read(dto.ReadContainersRequest) (dto.ReadContainersResponse, error)
}
