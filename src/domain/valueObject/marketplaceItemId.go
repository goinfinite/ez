package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type MarketplaceItemId uint16

func NewMarketplaceItemId(value interface{}) (
	marketplaceItemId MarketplaceItemId, err error,
) {
	uintValue, err := voHelper.InterfaceToUint16(value)
	if err != nil {
		return marketplaceItemId, errors.New("MarketplaceItemIdMustBeUint16")
	}

	return MarketplaceItemId(uintValue), nil
}

func (vo MarketplaceItemId) Uint16() uint16 {
	return uint16(vo)
}

func (vo MarketplaceItemId) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}
