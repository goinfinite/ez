package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type BackupTaskId uint64

func NewBackupTaskId(value interface{}) (BackupTaskId, error) {
	taskId, err := voHelper.InterfaceToUint64(value)
	if err != nil {
		return 0, errors.New("InvalidBackupTaskId")
	}

	return BackupTaskId(taskId), nil
}

func (vo BackupTaskId) Uint64() uint64 {
	return uint64(vo)
}

func (vo BackupTaskId) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}
