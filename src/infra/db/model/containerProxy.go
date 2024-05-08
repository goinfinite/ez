package dbModel

import "time"

type ContainerProxy struct {
	ID                   uint `gorm:"primarykey"`
	ContainerId          string
	ContainerHostname    string
	ContainerPrivatePort uint
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func (ContainerProxy) TableName() string {
	return "container_proxies"
}

func NewContainerProxy(
	id uint,
	containerId string,
	containerHostname string,
	containerPrivatePort uint,
) ContainerProxy {
	proxyModel := ContainerProxy{
		ContainerId:          containerId,
		ContainerHostname:    containerHostname,
		ContainerPrivatePort: containerPrivatePort,
	}

	if id != 0 {
		proxyModel.ID = id
	}

	return proxyModel
}
