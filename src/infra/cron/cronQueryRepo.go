package cronInfra

import (
	"errors"
	"time"

	"github.com/goinfinite/ez/src/domain/valueObject"
	"github.com/hashicorp/cronexpr"
)

type CronQueryRepo struct{}

func (repo *CronQueryRepo) ReadNextRun(
	cronSchedule valueObject.CronSchedule,
) (nextRunAt valueObject.UnixTime, err error) {
	cronExpression, err := cronexpr.Parse(cronSchedule.String())
	if err != nil {
		return nextRunAt, errors.New("ParseCronScheduleError: " + err.Error())
	}

	return valueObject.NewUnixTimeWithGoTime(
		cronExpression.Next(time.Now()),
	), nil
}
