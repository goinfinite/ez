package valueObject

import (
	"errors"
	"strings"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type BackupJobType string

const (
	BackupJobTypeFull        BackupJobType = "full"
	BackupJobTypeIncremental BackupJobType = "incremental"
)

func NewBackupJobType(value interface{}) (
	jobType BackupJobType, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return jobType, errors.New("BackupJobTypeMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	stringValueVo := BackupJobType(stringValue)
	switch stringValueVo {
	case BackupJobTypeFull, BackupJobTypeIncremental:
		return stringValueVo, nil
	default:
		return jobType, errors.New("InvalidBackupJobType")
	}
}

func (vo BackupJobType) String() string {
	return string(vo)
}
