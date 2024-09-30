package repository

import (
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type ContainerProxyCmdRepo interface {
	Create(containerId valueObject.ContainerId) error
	Delete(containerId valueObject.ContainerId) error
}
