package repository

import (
	"github.com/speedianet/sfm/src/domain/entity"
)

type ContainerQueryRepo interface {
	Get() ([]entity.Container, error)
}
