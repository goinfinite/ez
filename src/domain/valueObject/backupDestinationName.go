package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

const backupDestinationNameRegex string = `^\w[\w\ \-]{1,128}\w$`

type BackupDestinationName string

func NewBackupDestinationName(value interface{}) (name BackupDestinationName, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return name, errors.New("BackupDestinationNameMustBeString")
	}

	re := regexp.MustCompile(backupDestinationNameRegex)
	if !re.MatchString(stringValue) {
		return name, errors.New("InvalidBackupDestinationName")
	}

	return BackupDestinationName(stringValue), nil
}

func (vo BackupDestinationName) String() string {
	return string(vo)
}
