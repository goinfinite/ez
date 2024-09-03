package repository

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/valueObject"
)

type AuthQueryRepo interface {
	IsLoginValid(dto.CreateSessionToken) bool
	ReadAccessTokenDetails(valueObject.AccessTokenValue) (dto.AccessTokenDetails, error)
}
