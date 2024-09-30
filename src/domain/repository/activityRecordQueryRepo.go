package repository

import (
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
)

type ActivityRecordQueryRepo interface {
	Read(readDto dto.ReadActivityRecords) ([]entity.ActivityRecord, error)
}
