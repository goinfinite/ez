package repository

import (
	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

type AccCmdRepo interface {
	Add(addUser dto.AddUser) error
	Delete(userId valueObject.UserId) error
	UpdatePassword(userId valueObject.UserId, password valueObject.Password) error
	UpdateApiKey(userId valueObject.UserId) (valueObject.AccessTokenStr, error)
}
