package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

const scheduledTaskNameRegex string = `^[a-zA-Z][\w\-]{1,256}[\w\-\ ]{0,512}$`

type ScheduledTaskName string

func NewScheduledTaskName(value interface{}) (name ScheduledTaskName, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return name, errors.New("ScheduledTaskNameMustBeString")
	}

	re := regexp.MustCompile(scheduledTaskNameRegex)
	if !re.MatchString(stringValue) {
		return name, errors.New("InvalidScheduledTaskName")
	}

	return ScheduledTaskName(stringValue), nil
}

func (vo ScheduledTaskName) String() string {
	return string(vo)
}
