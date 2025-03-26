package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
)

var AccountsDefaultPagination dto.Pagination = dto.Pagination{
	PageNumber:   0,
	ItemsPerPage: 10,
}

func ReadAccounts(
	accountQueryRepo repository.AccountQueryRepo,
	requestDto dto.ReadAccountsRequest,
) (dto.ReadAccountsResponse, error) {
	responseDto, err := accountQueryRepo.Read(requestDto)
	if err != nil {
		slog.Error("ReadAccountsInfraError", slog.Any("error", err))
		return responseDto, errors.New("ReadAccountsInfraError")
	}

	return responseDto, nil
}
