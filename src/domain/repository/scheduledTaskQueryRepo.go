package repository

import (
	"github.com/speedianet/control/src/domain/entity"
)

type ScheduledTaskQueryRepo interface {
	Get() ([]entity.ScheduledTask, error)
}
