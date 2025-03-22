package valueObject

import (
	"errors"
	"strings"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type PaginationSortDirection string

const (
	PaginationSortDirectionAsc  PaginationSortDirection = "asc"
	PaginationSortDirectionDesc PaginationSortDirection = "desc"
)

func NewPaginationSortDirection(value interface{}) (
	sortDirection PaginationSortDirection, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return sortDirection, errors.New("PaginationSortDirectionMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	stringValueVo := PaginationSortDirection(stringValue)
	switch stringValueVo {
	case PaginationSortDirectionAsc, PaginationSortDirectionDesc:
		return stringValueVo, nil
	default:
		return sortDirection, errors.New("InvalidPaginationSortDirection")
	}
}

func (vo PaginationSortDirection) String() string {
	return string(vo)
}
