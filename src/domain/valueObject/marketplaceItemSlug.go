package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

const marketplaceItemSlugRegex string = `^[a-z0-9\_\-]{2,64}$`

type MarketplaceItemSlug string

func NewMarketplaceItemSlug(value interface{}) (slug MarketplaceItemSlug, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return slug, errors.New("MarketplaceItemSlugMustBeString")
	}

	re := regexp.MustCompile(marketplaceItemSlugRegex)
	if !re.MatchString(stringValue) {
		return slug, errors.New("InvalidMarketplaceItemSlug")
	}

	return MarketplaceItemSlug(stringValue), nil
}

func (vo MarketplaceItemSlug) String() string {
	return string(vo)
}
