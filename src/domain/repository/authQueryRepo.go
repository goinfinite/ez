package repository

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/valueObject"
)

type AuthQueryRepo interface {
	IsLoginValid(login dto.Login) bool
	ReadAccessTokenDetails(
		token valueObject.AccessTokenValue,
	) (dto.AccessTokenDetails, error)
}
