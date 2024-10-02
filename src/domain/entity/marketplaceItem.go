package entity

import "github.com/goinfinite/ez/src/domain/valueObject"

type MarketplaceItem struct {
	ReceiptVersion     valueObject.MarketplaceItemReceiptVersion `json:"receiptVersion"`
	Slug               valueObject.MarketplaceItemSlug           `json:"slug"`
	Name               valueObject.MarketplaceItemName           `json:"name"`
	Type               valueObject.MarketplaceItemType           `json:"type"`
	Description        valueObject.MarketplaceItemDescription    `json:"description"`
	LaunchScript       valueObject.LaunchScript                  `json:"launchScript"`
	EstimatedSizeBytes valueObject.Byte                          `json:"estimatedSizeBytes"`
	AvatarUrl          valueObject.Url                           `json:"avatarUrl,omitempty"`
}

func NewMarketplaceItem(
	receiptVersion valueObject.MarketplaceItemReceiptVersion,
	slug valueObject.MarketplaceItemSlug,
	name valueObject.MarketplaceItemName,
	itemType valueObject.MarketplaceItemType,
	description valueObject.MarketplaceItemDescription,
	launchScript valueObject.LaunchScript,
	estimatedSizeBytes valueObject.Byte,
	avatarUrl valueObject.Url,
) MarketplaceItem {
	return MarketplaceItem{
		ReceiptVersion:     receiptVersion,
		Slug:               slug,
		Name:               name,
		Type:               itemType,
		Description:        description,
		LaunchScript:       launchScript,
		EstimatedSizeBytes: estimatedSizeBytes,
		AvatarUrl:          avatarUrl,
	}
}
