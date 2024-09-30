package repository

import (
	"github.com/goinfinite/ez/src/domain/dto"
)

type ActivityRecordCmdRepo interface {
	Create(createDto dto.CreateActivityRecord) error
	Delete(deleteDto dto.DeleteActivityRecords) error
}
