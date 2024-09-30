package repository

import (
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type AuthQueryRepo interface {
	IsLoginValid(dto.CreateSessionToken) bool
	ReadAccessTokenDetails(valueObject.AccessTokenValue) (dto.AccessTokenDetails, error)
}
