package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

type ActivityRecordId uint64

func NewActivityRecordId(value interface{}) (ActivityRecordId, error) {
	id, err := voHelper.InterfaceToUint64(value)
	if err != nil {
		return 0, errors.New("InvalidActivityRecordId")
	}

	return ActivityRecordId(id), nil
}

func (vo ActivityRecordId) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}
