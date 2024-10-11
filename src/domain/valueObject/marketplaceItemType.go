package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type MarketplaceItemType string

var ValidMarketplaceItemTypes = []string{
	"app", "framework", "stack",
}

func NewMarketplaceItemType(value interface{}) (
	marketplaceItemType MarketplaceItemType, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return marketplaceItemType, errors.New("MarketplaceItemTypeMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	if !slices.Contains(ValidMarketplaceItemTypes, stringValue) {
		return marketplaceItemType, errors.New("InvalidMarketplaceItemType")
	}

	return MarketplaceItemType(stringValue), nil
}

func (vo MarketplaceItemType) String() string {
	return string(vo)
}
