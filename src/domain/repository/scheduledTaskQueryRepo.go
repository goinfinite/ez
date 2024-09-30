package repository

import (
	"github.com/goinfinite/ez/src/domain/dto"
)

type ScheduledTaskQueryRepo interface {
	Read(dto.ReadScheduledTasksRequest) (dto.ReadScheduledTasksResponse, error)
}
