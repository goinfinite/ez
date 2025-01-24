package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type BackupTaskArchiveId uint64

func NewBackupTaskArchiveId(value interface{}) (BackupTaskArchiveId, error) {
	archiveId, err := voHelper.InterfaceToUint64(value)
	if err != nil {
		return 0, errors.New("InvalidBackupTaskArchiveId")
	}

	return BackupTaskArchiveId(archiveId), nil
}

func (vo BackupTaskArchiveId) Uint64() uint64 {
	return uint64(vo)
}

func (vo BackupTaskArchiveId) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}
