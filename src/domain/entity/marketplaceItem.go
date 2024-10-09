package entity

import "github.com/goinfinite/ez/src/domain/valueObject"

type MarketplaceItem struct {
	ManifestVersion    valueObject.MarketplaceItemManifestVersion `json:"manifestVersion"`
	Slugs              []valueObject.MarketplaceItemSlug          `json:"slugs"`
	Name               valueObject.MarketplaceItemName            `json:"name"`
	Type               valueObject.MarketplaceItemType            `json:"type"`
	Description        valueObject.MarketplaceItemDescription     `json:"description"`
	LaunchScript       valueObject.LaunchScript                   `json:"launchScript"`
	EstimatedSizeBytes *valueObject.Byte                          `json:"estimatedSizeBytes,omitempty"`
	AvatarUrl          *valueObject.Url                           `json:"avatarUrl,omitempty"`
}

func NewMarketplaceItem(
	manifestVersion valueObject.MarketplaceItemManifestVersion,
	slugs []valueObject.MarketplaceItemSlug,
	name valueObject.MarketplaceItemName,
	itemType valueObject.MarketplaceItemType,
	description valueObject.MarketplaceItemDescription,
	launchScript valueObject.LaunchScript,
	estimatedSizeBytes *valueObject.Byte,
	avatarUrl *valueObject.Url,
) MarketplaceItem {
	return MarketplaceItem{
		ManifestVersion:    manifestVersion,
		Slugs:              slugs,
		Name:               name,
		Type:               itemType,
		Description:        description,
		LaunchScript:       launchScript,
		EstimatedSizeBytes: estimatedSizeBytes,
		AvatarUrl:          avatarUrl,
	}
}
