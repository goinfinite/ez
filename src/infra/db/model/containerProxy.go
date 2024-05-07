package dbModel

import "time"

type ContainerProxy struct {
	ID                uint `gorm:"primarykey"`
	ContainerId       string
	ContainerHostname string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (ContainerProxy) TableName() string {
	return "container_proxies"
}

func NewContainerProxy(
	id uint,
	containerId string,
	containerHostname string,
) ContainerProxy {
	proxyModel := ContainerProxy{
		ContainerId:       containerId,
		ContainerHostname: containerHostname,
	}

	if id != 0 {
		proxyModel.ID = id
	}

	return proxyModel
}
