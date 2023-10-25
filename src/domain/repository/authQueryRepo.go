package repository

import (
	"github.com/goinfinite/fleet/src/domain/dto"
	"github.com/goinfinite/fleet/src/domain/valueObject"
)

type AuthQueryRepo interface {
	IsLoginValid(login dto.Login) bool
	GetAccessTokenDetails(
		token valueObject.AccessTokenStr,
	) (dto.AccessTokenDetails, error)
}
