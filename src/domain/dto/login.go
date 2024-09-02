package dto

import "github.com/speedianet/control/src/domain/valueObject"

type Login struct {
	Username          valueObject.Username   `json:"username"`
	Password          valueObject.Password   `json:"password"`
	OperatorIpAddress *valueObject.IpAddress `json:"-"`
}

func NewLogin(
	username valueObject.Username,
	password valueObject.Password,
	operatorIpAddress *valueObject.IpAddress,
) Login {
	return Login{
		Username:          username,
		Password:          password,
		OperatorIpAddress: operatorIpAddress,
	}
}
