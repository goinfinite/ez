package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
)

func UpdateAccount(
	accountQueryRepo repository.AccountQueryRepo,
	accountCmdRepo repository.AccountCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	updateDto dto.UpdateAccount,
) error {
	_, err := accountQueryRepo.ReadFirst(dto.ReadAccountsRequest{
		AccountId: &updateDto.AccountId,
	})
	if err != nil {
		return errors.New("AccountNotFound")
	}

	if updateDto.Password != nil {
		err = accountCmdRepo.UpdatePassword(updateDto.AccountId, *updateDto.Password)
		if err != nil {
			slog.Error(
				"UpdateAccountPasswordInfraError",
				slog.String("accountId", updateDto.AccountId.String()),
				slog.Any("error", err),
			)
			return errors.New("UpdateAccountPasswordInfraError")
		}

		NewCreateSecurityActivityRecord(activityRecordCmdRepo).UpdateAccount(updateDto)
	}

	if updateDto.Quota == nil {
		return nil
	}

	err = accountCmdRepo.UpdateQuota(updateDto.AccountId, *updateDto.Quota)
	if err != nil {
		slog.Error(
			"UpdateAccountQuotaInfraError",
			slog.String("accountId", updateDto.AccountId.String()),
			slog.Any("error", err),
		)
		return errors.New("UpdateAccountQuotaInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).UpdateAccount(updateDto)
	return nil
}
