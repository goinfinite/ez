package dto

import "github.com/goinfinite/ez/src/domain/valueObject"

type Pagination struct {
	PageNumber    uint32                               `json:"pageNumber"`
	ItemsPerPage  uint16                               `json:"itemsPerPage"`
	SortBy        *valueObject.PaginationSortBy        `json:"sortBy"`
	SortDirection *valueObject.PaginationSortDirection `json:"sortDirection"`
	LastSeenId    *valueObject.PaginationLastSeenId    `json:"lastSeenId"`
	PagesTotal    *uint32                              `json:"pagesTotal"`
	ItemsTotal    *uint64                              `json:"itemsTotal"`
}

var PaginationSingleItem = Pagination{
	PageNumber:   0,
	ItemsPerPage: 1,
}

var PaginationUnpaginated = Pagination{
	PageNumber:   0,
	ItemsPerPage: 1000,
}
