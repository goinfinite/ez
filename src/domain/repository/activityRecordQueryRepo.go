package repository

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
)

type ActivityRecordQueryRepo interface {
	Read(readDto dto.ReadActivityRecords) ([]entity.ActivityRecord, error)
}
