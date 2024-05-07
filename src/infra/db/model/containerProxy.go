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
