package entity

import "github.com/goinfinite/ez/src/domain/valueObject"

type MarketplaceItem struct {
	ManifestVersion      valueObject.MarketplaceItemManifestVersion `json:"manifestVersion"`
	Slugs                []valueObject.MarketplaceItemSlug          `json:"slugs"`
	Name                 valueObject.MarketplaceItemName            `json:"name"`
	Type                 valueObject.MarketplaceItemType            `json:"type"`
	Description          valueObject.MarketplaceItemDescription     `json:"description"`
	RegistryImageAddress valueObject.ContainerImageAddress          `json:"registryImageAddress"`
	LaunchScript         valueObject.LaunchScript                   `json:"launchScript"`
	MinimumCpuMillicores *valueObject.Millicores                    `json:"minimumCpuMillicores,omitempty"`
	MinimumMemoryBytes   *valueObject.Byte                          `json:"minimumMemoryBytes,omitempty"`
	EstimatedSizeBytes   *valueObject.Byte                          `json:"estimatedSizeBytes,omitempty"`
	AvatarUrl            *valueObject.Url                           `json:"avatarUrl,omitempty"`
}

func NewMarketplaceItem(
	manifestVersion valueObject.MarketplaceItemManifestVersion,
	slugs []valueObject.MarketplaceItemSlug,
	name valueObject.MarketplaceItemName,
	itemType valueObject.MarketplaceItemType,
	description valueObject.MarketplaceItemDescription,
	registryImageAddress valueObject.ContainerImageAddress,
	launchScript valueObject.LaunchScript,
	minimumCpuMillicores *valueObject.Millicores,
	minimumMemoryBytes *valueObject.Byte,
	estimatedSizeBytes *valueObject.Byte,
	avatarUrl *valueObject.Url,
) MarketplaceItem {
	return MarketplaceItem{
		ManifestVersion:      manifestVersion,
		Slugs:                slugs,
		Name:                 name,
		Type:                 itemType,
		Description:          description,
		RegistryImageAddress: registryImageAddress,
		LaunchScript:         launchScript,
		MinimumCpuMillicores: minimumCpuMillicores,
		MinimumMemoryBytes:   minimumMemoryBytes,
		EstimatedSizeBytes:   estimatedSizeBytes,
		AvatarUrl:            avatarUrl,
	}
}
