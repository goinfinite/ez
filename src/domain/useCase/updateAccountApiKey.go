package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

func UpdateAccountApiKey(
	accountQueryRepo repository.AccountQueryRepo,
	accountCmdRepo repository.AccountCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	updateDto dto.UpdateAccount,
) (newKey valueObject.AccessTokenValue, err error) {
	_, err = accountQueryRepo.ReadFirst(dto.ReadAccountsRequest{
		AccountId: &updateDto.AccountId,
	})
	if err != nil {
		return newKey, errors.New("AccountNotFound")
	}

	newKey, err = accountCmdRepo.UpdateApiKey(updateDto.AccountId)
	if err != nil {
		slog.Error(
			"UpdateAccountApiKeyInfraError",
			slog.String("accountId", updateDto.AccountId.String()),
			slog.Any("error", err),
		)
		return newKey, errors.New("UpdateAccountApiKeyInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).UpdateAccount(updateDto)

	return newKey, nil
}
