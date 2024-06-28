package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

type SecurityEventId uint

func NewSecurityEventId(value interface{}) (SecurityEventId, error) {
	id, err := voHelper.InterfaceToUint(value)
	if err != nil {
		return 0, errors.New("InvalidSecurityEventId")
	}

	return SecurityEventId(id), nil
}

func (vo SecurityEventId) Get() uint64 {
	return uint64(vo)
}

func (vo SecurityEventId) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}
