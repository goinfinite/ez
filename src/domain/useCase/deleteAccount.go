package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
)

func DeleteAccount(
	accountQueryRepo repository.AccountQueryRepo,
	accountCmdRepo repository.AccountCmdRepo,
	containerQueryRepo repository.ContainerQueryRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteAccount,
) error {
	_, err := accountQueryRepo.ReadById(deleteDto.AccountId)
	if err != nil {
		return errors.New("AccountNotFound")
	}

	containers, err := containerQueryRepo.ReadByAccountId(deleteDto.AccountId)
	if err != nil {
		slog.Error("ReadContainersByAccountIdInfraError", slog.Any("error", err))
		return errors.New("ReadContainersByAccountIdInfraError")
	}

	if len(containers) > 0 {
		return errors.New("AccountHasContainers")
	}

	err = accountCmdRepo.Delete(deleteDto.AccountId)
	if err != nil {
		slog.Error("DeleteAccountInfraError", slog.Any("error", err))
		return errors.New("DeleteAccountInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).DeleteAccount(deleteDto)

	return nil
}
