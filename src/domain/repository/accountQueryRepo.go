package repository

import (
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type AccountQueryRepo interface {
	Read(dto.ReadAccountsRequest) (dto.ReadAccountsResponse, error)
	ReadByUsername(username valueObject.UnixUsername) (entity.Account, error)
	ReadById(accountId valueObject.AccountId) (entity.Account, error)
}
