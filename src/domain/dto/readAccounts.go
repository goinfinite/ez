package dto

import (
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type ReadAccountsRequest struct {
	Pagination      Pagination                `json:"pagination"`
	AccountId       *valueObject.AccountId    `json:"id,omitempty"`
	AccountUsername *valueObject.UnixUsername `json:"username,omitempty"`
}

type ReadAccountsResponse struct {
	Pagination Pagination       `json:"pagination"`
	Accounts   []entity.Account `json:"accounts"`
}
