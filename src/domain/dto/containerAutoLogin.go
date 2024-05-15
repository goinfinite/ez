package dto

import "github.com/speedianet/control/src/domain/valueObject"

type ContainerAutoLogin struct {
	ContainerId valueObject.ContainerId `json:"containerId"`
	IpAddress   valueObject.IpAddress   `json:"-"`
}

func NewContainerAutoLogin(
	containerId valueObject.ContainerId,
	ipAddress valueObject.IpAddress,
) ContainerAutoLogin {
	return ContainerAutoLogin{
		ContainerId: containerId,
		IpAddress:   ipAddress,
	}
}
