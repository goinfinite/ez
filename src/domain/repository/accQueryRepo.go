package repository

import (
	"github.com/speedianet/sfm/src/domain/entity"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

type AccQueryRepo interface {
	Get() ([]entity.AccountDetails, error)
	GetByUsername(
		username valueObject.Username,
	) (entity.AccountDetails, error)
	GetById(
		userId valueObject.UserId,
	) (entity.AccountDetails, error)
}
