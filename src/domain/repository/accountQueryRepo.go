package repository

import (
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
)

type AccountQueryRepo interface {
	Read(dto.ReadAccountsRequest) (dto.ReadAccountsResponse, error)
	ReadFirst(dto.ReadAccountsRequest) (entity.Account, error)
}
