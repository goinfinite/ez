package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/sfm/src/domain/valueObject/helper"
)

type ResourceProfileId uint64

func NewResourceProfileId(value interface{}) (ResourceProfileId, error) {
	rpId, err := voHelper.InterfaceToUint(value)
	if err != nil {
		return 0, errors.New("InvalidResourceProfileId")
	}

	return ResourceProfileId(rpId), nil
}

func NewResourceProfileIdPanic(value interface{}) ResourceProfileId {
	rpId, err := NewResourceProfileId(value)
	if err != nil {
		panic(err)
	}
	return rpId
}

func (id ResourceProfileId) Get() uint64 {
	return uint64(id)
}

func (id ResourceProfileId) String() string {
	return strconv.FormatUint(uint64(id), 10)
}
