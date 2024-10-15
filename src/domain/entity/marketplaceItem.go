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
	MinimumMemoryBytes   *valueObject.Byte                          `json:"minimumMemoryBytes,omitempty"`
	MinimumCpuMillicores *valueObject.Millicores                    `json:"minimumCpuMillicores,omitempty"`
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
	minimumMemoryBytes *valueObject.Byte,
	minimumCpuMillicores *valueObject.Millicores,
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
		MinimumMemoryBytes:   minimumMemoryBytes,
		MinimumCpuMillicores: minimumCpuMillicores,
		EstimatedSizeBytes:   estimatedSizeBytes,
		AvatarUrl:            avatarUrl,
	}
}
