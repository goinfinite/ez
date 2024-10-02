package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type MarketplaceItemReceiptVersion string

var validMarketplaceItemReceiptVersions = []string{
	"v1",
}

func NewMarketplaceItemReceiptVersion(value interface{}) (
	version MarketplaceItemReceiptVersion, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return version, errors.New("MarketplaceItemReceiptVersionMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	if !slices.Contains(validMarketplaceItemReceiptVersions, stringValue) {
		return version, errors.New("InvalidMarketplaceItemReceiptVersion")
	}

	return MarketplaceItemReceiptVersion(stringValue), nil
}

func (vo MarketplaceItemReceiptVersion) String() string {
	return string(vo)
}
