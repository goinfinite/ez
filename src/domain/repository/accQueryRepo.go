package repository

import (
	"github.com/goinfinite/fleet/src/domain/entity"
	"github.com/goinfinite/fleet/src/domain/valueObject"
)

type AccQueryRepo interface {
	Get() ([]entity.Account, error)
	GetByUsername(
		username valueObject.Username,
	) (entity.Account, error)
	GetById(
		accountId valueObject.AccountId,
	) (entity.Account, error)
}
