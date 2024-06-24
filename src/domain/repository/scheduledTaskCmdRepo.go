package repository

import (
	"github.com/speedianet/control/src/domain/dto"
)

type ScheduledTaskCmdRepo interface {
	Create(createDto dto.CreateScheduledTask) error
}
