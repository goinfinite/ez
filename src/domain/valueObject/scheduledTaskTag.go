package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

const scheduledTaskTagRegex string = `^[a-zA-Z][\w\-]{1,256}$`

type ScheduledTaskTag string

func NewScheduledTaskTag(value interface{}) (ScheduledTaskTag, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("ScheduledTaskTagMustBeString")
	}

	re := regexp.MustCompile(scheduledTaskTagRegex)
	if !re.MatchString(stringValue) {
		return "", errors.New("InvalidScheduledTaskTag")
	}

	return ScheduledTaskTag(stringValue), nil
}

func (vo ScheduledTaskTag) String() string {
	return string(vo)
}
