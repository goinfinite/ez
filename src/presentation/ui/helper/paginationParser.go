package uiHelper

import (
	"github.com/goinfinite/ez/src/domain/valueObject"
	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
	"github.com/labstack/echo/v4"
)

func PaginationParser(
	echoContext echo.Context,
	itemName string,
	defaultSortColumnName string,
) map[string]interface{} {
	itemNamePageNumber := uint16(0)
	if echoContext.QueryParam(itemName+"PageNumber") != "" {
		itemNamePageNumber, _ = voHelper.InterfaceToUint16(
			echoContext.QueryParam(itemName + "PageNumber"),
		)
	}

	itemNameItemsPerPage := uint16(10)
	if echoContext.QueryParam(itemName+"ItemsPerPage") != "" {
		itemNameItemsPerPage, _ = voHelper.InterfaceToUint16(
			echoContext.QueryParam(itemName + "ItemsPerPage"),
		)
	}

	itemNameSortByStr := defaultSortColumnName
	if echoContext.QueryParam(itemName+"SortBy") != "" {
		itemNameSortBy, err := valueObject.NewPaginationSortBy(
			echoContext.QueryParam(itemName + "SortBy"),
		)
		if err == nil {
			itemNameSortByStr = itemNameSortBy.String()
		}
	}

	itemNameSortDirectionStr := "asc"
	if echoContext.QueryParam(itemName+"SortDirection") != "" {
		itemNameSortDirection, err := valueObject.NewPaginationSortDirection(
			echoContext.QueryParam(itemName + "SortDirection"),
		)
		if err == nil {
			itemNameSortDirectionStr = itemNameSortDirection.String()
		}
	}

	return map[string]interface{}{
		"pageNumber":    itemNamePageNumber,
		"itemsPerPage":  itemNameItemsPerPage,
		"sortBy":        itemNameSortByStr,
		"sortDirection": itemNameSortDirectionStr,
	}
}
