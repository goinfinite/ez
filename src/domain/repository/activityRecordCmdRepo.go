package repository

import (
	"github.com/speedianet/control/src/domain/dto"
)

type ActivityRecordCmdRepo interface {
	Create(createDto dto.CreateActivityRecord) error
	Delete(deleteDto dto.DeleteActivityRecords) error
}
