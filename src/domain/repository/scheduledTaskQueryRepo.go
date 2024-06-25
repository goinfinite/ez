package repository

import (
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type ScheduledTaskQueryRepo interface {
	Get() ([]entity.ScheduledTask, error)
	GetById(id valueObject.ScheduledTaskId) (entity.ScheduledTask, error)
}
