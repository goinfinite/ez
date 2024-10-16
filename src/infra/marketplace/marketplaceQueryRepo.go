package marketplaceInfra

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
	infraEnvs "github.com/goinfinite/ez/src/infra/envs"
	infraHelper "github.com/goinfinite/ez/src/infra/helper"
	"gopkg.in/yaml.v3"
)

type MarketplaceQueryRepo struct {
}

func NewMarketplaceQueryRepo() *MarketplaceQueryRepo {
	return &MarketplaceQueryRepo{}
}

func (repo *MarketplaceQueryRepo) readFileContentToMap(
	filePath valueObject.UnixFilePath,
) (fileContentMap map[string]interface{}, err error) {
	itemFileExt, err := filePath.ReadFileExtension()
	if err != nil {
		return fileContentMap, err
	}

	supportedFileExtensions := []string{"yml", "yaml", "json"}
	if !slices.Contains(supportedFileExtensions, itemFileExt.String()) {
		return fileContentMap, errors.New("UnsupportedMarketplaceItemFileExtension")
	}

	fileHandler, err := os.Open(filePath.String())
	if err != nil {
		return fileContentMap, errors.New("OpenMarketplaceItemFileError: " + err.Error())
	}

	isYamlFile := itemFileExt == "yml" || itemFileExt == "yaml"
	if isYamlFile {
		itemYamlDecoder := yaml.NewDecoder(fileHandler)
		err = itemYamlDecoder.Decode(&fileContentMap)
		if err != nil {
			return fileContentMap, errors.New("DecodeMarketplaceItemYamlError: " + err.Error())
		}

		return fileContentMap, nil
	}

	itemJsonDecoder := json.NewDecoder(fileHandler)
	err = itemJsonDecoder.Decode(&fileContentMap)
	if err != nil {
		return fileContentMap, errors.New("DecodeMarketplaceItemJsonError: " + err.Error())
	}

	return fileContentMap, nil
}

func (repo *MarketplaceQueryRepo) itemFactory(
	itemFilePath valueObject.UnixFilePath,
) (itemEntity entity.MarketplaceItem, err error) {
	itemMap, err := repo.readFileContentToMap(itemFilePath)
	if err != nil {
		return itemEntity, err
	}

	requiredFields := []string{
		"manifestVersion", "slugs", "name", "type",
		"description", "registryImageAddress", "launchScript",
	}
	missingFields := []string{}
	for _, requiredField := range requiredFields {
		if _, exists := itemMap[requiredField]; !exists {
			missingFields = append(missingFields, requiredField)
		}
	}
	if len(missingFields) > 0 {
		return itemEntity, errors.New("MissingItemFields: " + strings.Join(missingFields, ", "))
	}

	manifestVersion, err := valueObject.NewMarketplaceItemManifestVersion(
		itemMap["manifestVersion"],
	)
	if err != nil {
		return itemEntity, err
	}

	itemSlugs := []valueObject.MarketplaceItemSlug{}
	if itemMap["slugs"] != nil {
		rawItemSlugs, assertOk := itemMap["slugs"].([]interface{})
		if !assertOk {
			return itemEntity, errors.New("InvalidMarketplaceItemSlugs")
		}
		for _, rawItemSlug := range rawItemSlugs {
			itemSlug, err := valueObject.NewMarketplaceItemSlug(rawItemSlug)
			if err != nil {
				slog.Debug(err.Error(), slog.Any("rawItemSlug", rawItemSlug))
				continue
			}
			itemSlugs = append(itemSlugs, itemSlug)
		}
	}

	itemName, err := valueObject.NewMarketplaceItemName(itemMap["name"])
	if err != nil {
		return itemEntity, err
	}

	itemType, err := valueObject.NewMarketplaceItemType(itemMap["type"])
	if err != nil {
		return itemEntity, err
	}

	itemDescription, err := valueObject.NewMarketplaceItemDescription(
		itemMap["description"],
	)
	if err != nil {
		return itemEntity, err
	}

	registryImageAddress, err := valueObject.NewContainerImageAddress(
		itemMap["registryImageAddress"],
	)
	if err != nil {
		return itemEntity, err
	}

	rawLaunchScriptSlice, assertOk := itemMap["launchScript"].([]interface{})
	if !assertOk {
		rawLaunchScriptStr, assertOk := itemMap["launchScript"].(string)
		if !assertOk {
			return itemEntity, errors.New("InvalidMarketplaceItemLaunchScript")
		}
		rawLaunchScriptSlice = []interface{}{rawLaunchScriptStr}
	}

	rawLaunchScript := ""
	for _, rawLaunchScriptStep := range rawLaunchScriptSlice {
		rawLaunchScriptStep, assertOk := rawLaunchScriptStep.(string)
		if !assertOk {
			slog.Debug(
				"InvalidMarketplaceItemLaunchScriptStep",
				slog.Any("rawLaunchScriptStep", rawLaunchScriptStep),
			)
			return itemEntity, errors.New("InvalidMarketplaceItemLaunchScript")
		}
		rawLaunchScript += rawLaunchScriptStep + "\n"
	}

	launchScript, err := valueObject.NewLaunchScript(rawLaunchScript)
	if err != nil {
		return itemEntity, err
	}

	var minimumCpuMillicoresPtr *valueObject.Millicores
	if itemMap["minimumCpuMillicores"] != nil {
		minimumCpuMillicores, err := valueObject.NewMillicores(itemMap["minimumCpuMillicores"])
		if err != nil {
			return itemEntity, err
		}
		minimumCpuMillicoresPtr = &minimumCpuMillicores
	}

	var minimumMemoryBytesPtr *valueObject.Byte
	if itemMap["minimumMemoryBytes"] != nil {
		minimumMemoryBytes, err := valueObject.NewByte(itemMap["minimumMemoryBytes"])
		if err != nil {
			return itemEntity, err
		}
		minimumMemoryBytesPtr = &minimumMemoryBytes
	}

	var estimatedSizeBytesPtr *valueObject.Byte
	if itemMap["estimatedSizeBytes"] != nil {
		estimatedSizeBytes, err := valueObject.NewByte(itemMap["estimatedSizeBytes"])
		if err != nil {
			return itemEntity, err
		}
		estimatedSizeBytesPtr = &estimatedSizeBytes
	}

	var itemAvatarUrlPtr *valueObject.Url
	if itemMap["avatarUrl"] != nil {
		itemAvatarUrl, err := valueObject.NewUrl(itemMap["avatarUrl"])
		if err != nil {
			return itemEntity, err
		}
		itemAvatarUrlPtr = &itemAvatarUrl
	}

	return entity.NewMarketplaceItem(
		manifestVersion, itemSlugs, itemName, itemType, itemDescription, registryImageAddress,
		launchScript, minimumCpuMillicoresPtr, minimumMemoryBytesPtr, estimatedSizeBytesPtr,
		itemAvatarUrlPtr,
	), nil
}

func (repo *MarketplaceQueryRepo) Read(
	readDto dto.ReadMarketplaceItemsRequest,
) (responseDto dto.ReadMarketplaceItemsResponse, err error) {
	_, err = os.Stat(infraEnvs.MarketplaceDir)
	if err != nil {
		marketplaceCmdRepo := NewMarketplaceCmdRepo()
		err = marketplaceCmdRepo.Refresh()
		if err != nil {
			return responseDto, errors.New("RefreshMarketplaceError: " + err.Error())
		}
	}

	rawFilesList, err := infraHelper.RunCmdWithSubShell(
		"find " + infraEnvs.MarketplaceDir + " -type f " +
			"\\( -name '*.json' -o -name '*.yaml' -o -name '*.yml' \\) " +
			"-not -path '*/.*' -not -name '.*'",
	)
	if err != nil {
		return responseDto, errors.New("ReadMarketplaceFilesError: " + err.Error())
	}

	if len(rawFilesList) == 0 {
		return responseDto, errors.New("NoMarketplaceFilesFound")
	}

	rawFilesListParts := strings.Split(rawFilesList, "\n")
	if len(rawFilesListParts) == 0 {
		return responseDto, errors.New("NoMarketplaceFilesFound")
	}

	nothingToFilter := readDto.ItemSlug == nil &&
		readDto.ItemName == nil &&
		readDto.ItemType == nil

	itemsList := []entity.MarketplaceItem{}
	for _, rawFilePath := range rawFilesListParts {
		itemFilePath, err := valueObject.NewUnixFilePath(rawFilePath)
		if err != nil {
			slog.Error(err.Error(), slog.String("filePath", rawFilePath))
			continue
		}

		marketplaceItem, err := repo.itemFactory(itemFilePath)
		if err != nil {
			slog.Error(
				"MarketplaceItemFactoryError",
				slog.String("filePath", itemFilePath.String()), slog.Any("err", err),
			)
			continue
		}

		if len(itemsList) >= int(readDto.Pagination.ItemsPerPage) {
			break
		}

		if nothingToFilter {
			itemsList = append(itemsList, marketplaceItem)
			continue
		}

		if readDto.ItemSlug != nil {
			if !slices.Contains(marketplaceItem.Slugs, *readDto.ItemSlug) {
				continue
			}
		}

		if readDto.ItemName != nil && marketplaceItem.Name != *readDto.ItemName {
			continue
		}

		if readDto.ItemType != nil && marketplaceItem.Type != *readDto.ItemType {
			continue
		}

		itemsList = append(itemsList, marketplaceItem)
	}

	sortDirectionStr := "asc"
	if readDto.Pagination.SortDirection != nil {
		sortDirectionStr = readDto.Pagination.SortDirection.String()
	}

	if readDto.Pagination.SortBy != nil {
		slices.SortStableFunc(itemsList, func(a, b entity.MarketplaceItem) int {
			switch readDto.Pagination.SortBy.String() {
			case "name":
				if sortDirectionStr != "asc" {
					return strings.Compare(b.Name.String(), a.Name.String())
				}
				return strings.Compare(a.Name.String(), b.Name.String())
			case "type":
				if sortDirectionStr != "asc" {
					return strings.Compare(b.Type.String(), a.Type.String())
				}
				return strings.Compare(a.Type.String(), b.Type.String())
			default:
				return 0
			}
		})
	}

	itemsTotal := uint64(len(itemsList))
	pagesTotal := uint32(itemsTotal / uint64(readDto.Pagination.ItemsPerPage))

	paginationDto := readDto.Pagination
	paginationDto.ItemsTotal = &itemsTotal
	paginationDto.PagesTotal = &pagesTotal

	return dto.ReadMarketplaceItemsResponse{
		Pagination: paginationDto,
		Items:      itemsList,
	}, nil
}
