package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type BackupDestinationId uint64

func NewBackupDestinationId(value interface{}) (BackupDestinationId, error) {
	destinationId, err := voHelper.InterfaceToUint64(value)
	if err != nil {
		return 0, errors.New("InvalidBackupDestinationId")
	}

	return BackupDestinationId(destinationId), nil
}

func (vo BackupDestinationId) Uint64() uint64 {
	return uint64(vo)
}

func (vo BackupDestinationId) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}
