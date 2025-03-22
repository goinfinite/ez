package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type BackupJobId uint64

func NewBackupJobId(value interface{}) (BackupJobId, error) {
	jobId, err := voHelper.InterfaceToUint64(value)
	if err != nil {
		return 0, errors.New("InvalidBackupJobId")
	}

	return BackupJobId(jobId), nil
}

func (vo BackupJobId) Uint64() uint64 {
	return uint64(vo)
}

func (vo BackupJobId) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}
