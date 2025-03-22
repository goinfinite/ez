package dto

import "github.com/goinfinite/ez/src/domain/valueObject"

type CreateSessionToken struct {
	Username          valueObject.UnixUsername `json:"username"`
	Password          valueObject.Password     `json:"password"`
	OperatorIpAddress valueObject.IpAddress    `json:"-"`
}

func NewCreateSessionToken(
	username valueObject.UnixUsername,
	password valueObject.Password,
	operatorIpAddress valueObject.IpAddress,
) CreateSessionToken {
	return CreateSessionToken{
		Username:          username,
		Password:          password,
		OperatorIpAddress: operatorIpAddress,
	}
}
