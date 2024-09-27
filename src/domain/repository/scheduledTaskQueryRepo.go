package repository

import (
	"github.com/speedianet/control/src/domain/dto"
)

type ScheduledTaskQueryRepo interface {
	Read(dto.ReadScheduledTasksRequest) (dto.ReadScheduledTasksResponse, error)
}
