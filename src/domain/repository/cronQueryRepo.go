package repository

import (
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type CronQueryRepo interface {
	ReadNextRun(valueObject.CronSchedule) (valueObject.UnixTime, error)
}
