package marketplaceInfra

import (
	"encoding/json"
	"errors"
	"log/slog"
	"math/rand"
	"os"
	"slices"
	"strconv"
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

func (repo *MarketplaceQueryRepo) launchScriptFactory(
	rawLaunchScriptSlice []interface{},
) (launchScript valueObject.LaunchScript, err error) {
	rawLaunchScript := ""
	for stepIndex, rawLaunchScriptStep := range rawLaunchScriptSlice {
		rawLaunchScriptStep, assertOk := rawLaunchScriptStep.(string)
		if !assertOk {
			return launchScript, errors.New(
				"[" + strconv.Itoa(stepIndex) + "] InvalidMarketplaceItemLaunchScript",
			)
		}
		rawLaunchScript += rawLaunchScriptStep + "\n"
	}

	if len(rawLaunchScript) == 0 {
		return launchScript, errors.New("EmptyMarketplaceItemLaunchScript")
	}

	randomUsernames := []string{
		"spock", "kirk", "mccoy", "scotty", "uhura", "sulu", "chekov",
	}
	randomUsername := randomUsernames[rand.Intn(len(randomUsernames))]
	rawLaunchScript = strings.ReplaceAll(
		rawLaunchScript, "%randomUsername%", randomUsername,
	)
	rawLaunchScript = strings.ReplaceAll(
		rawLaunchScript, "%randomPassword%", infraHelper.GenPass(16),
	)
	rawLaunchScript = strings.ReplaceAll(
		rawLaunchScript, "%randomMail%", randomUsername+"@ufp.gov",
	)

	return valueObject.NewLaunchScript(rawLaunchScript)
}

func (repo *MarketplaceQueryRepo) itemFactory(
	itemFilePath valueObject.UnixFilePath,
) (itemEntity entity.MarketplaceItem, err error) {
	itemMap, err := repo.readFileContentToMap(itemFilePath)
	if err != nil {
		return itemEntity, err
	}

	requiredFields := []string{
		"manifestVersion", "slugs", "name", "type", "description", "registryImageAddress",
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

	itemId, _ := valueObject.NewMarketplaceItemId(0)
	if itemMap["id"] != nil {
		itemId, err = valueObject.NewMarketplaceItemId(itemMap["id"])
		if err != nil {
			return itemEntity, err
		}
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

	var launchScriptPtr *valueObject.LaunchScript
	if itemMap["launchScript"] != nil {
		rawLaunchScriptSlice, assertOk := itemMap["launchScript"].([]interface{})
		if !assertOk {
			rawLaunchScriptStr, assertOk := itemMap["launchScript"].(string)
			if !assertOk {
				return itemEntity, errors.New("InvalidMarketplaceItemLaunchScript")
			}
			rawLaunchScriptSlice = []interface{}{rawLaunchScriptStr}
		}

		launchScript, err := repo.launchScriptFactory(rawLaunchScriptSlice)
		if err != nil {
			return itemEntity, err
		}

		launchScriptPtr = &launchScript
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
		manifestVersion, itemId, itemSlugs, itemName, itemType, itemDescription,
		registryImageAddress, launchScriptPtr, minimumCpuMillicoresPtr, minimumMemoryBytesPtr,
		estimatedSizeBytesPtr, itemAvatarUrlPtr,
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

	itemsList := []entity.MarketplaceItem{}
	itemsIdsMap := map[uint16]struct{}{}
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

		itemIdUint16 := marketplaceItem.Id.Uint16()
		_, idAlreadyUsed := itemsIdsMap[itemIdUint16]
		if idAlreadyUsed {
			marketplaceItem.Id, _ = valueObject.NewMarketplaceItemId(0)
		}

		if len(itemsList) >= int(readDto.Pagination.ItemsPerPage) {
			break
		}

		if readDto.ItemId != nil && marketplaceItem.Id != *readDto.ItemId {
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

		if marketplaceItem.Id.Uint16() != 0 {
			itemsIdsMap[itemIdUint16] = struct{}{}
		}
	}

	itemsIdsSlice := []uint16{}
	for itemId := range itemsIdsMap {
		itemsIdsSlice = append(itemsIdsSlice, itemId)
	}
	slices.Sort(itemsIdsSlice)

	if len(itemsIdsSlice) == 0 {
		itemsIdsSlice = append(itemsIdsSlice, 0)
	}

	for itemIndex, marketplaceItem := range itemsList {
		if marketplaceItem.Id.Uint16() != 0 {
			continue
		}

		lastIdUsed := itemsIdsSlice[len(itemsIdsSlice)-1]
		nextAvailableId, err := valueObject.NewMarketplaceItemId(lastIdUsed + 1)
		if err != nil {
			slog.Error(
				"CreateNewMarketplaceItemIdError",
				slog.String("itemName", marketplaceItem.Name.String()),
				slog.Any("err", err),
			)
			continue
		}

		itemsList[itemIndex].Id = nextAvailableId
		itemsIdsSlice = append(itemsIdsSlice, nextAvailableId.Uint16())
	}

	sortDirectionStr := "asc"
	if readDto.Pagination.SortDirection != nil {
		sortDirectionStr = readDto.Pagination.SortDirection.String()
	}

	if readDto.Pagination.SortBy != nil {
		slices.SortStableFunc(itemsList, func(a, b entity.MarketplaceItem) int {
			firstElement := a
			secondElement := b
			if sortDirectionStr != "asc" {
				firstElement = b
				secondElement = a
			}

			switch readDto.Pagination.SortBy.String() {
			case "id":
				if firstElement.Id.Uint16() < secondElement.Id.Uint16() {
					return -1
				}
				if firstElement.Id.Uint16() > secondElement.Id.Uint16() {
					return 1
				}
				return 0
			case "name":
				return strings.Compare(firstElement.Name.String(), secondElement.Name.String())
			case "type":
				return strings.Compare(firstElement.Type.String(), secondElement.Type.String())
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
