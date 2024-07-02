package repository

import (
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type ScheduledTaskQueryRepo interface {
	Read() ([]entity.ScheduledTask, error)
	ReadById(id valueObject.ScheduledTaskId) (entity.ScheduledTask, error)
	ReadByStatus(status valueObject.ScheduledTaskStatus) ([]entity.ScheduledTask, error)
}
