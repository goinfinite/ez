package useCase

import (
	"errors"
	"log/slog"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func CreateAccount(
	accountQueryRepo repository.AccountQueryRepo,
	accountCmdRepo repository.AccountCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateAccount,
) error {
	_, err := accountQueryRepo.ReadByUsername(createDto.Username)
	if err == nil {
		return errors.New("AccountAlreadyExists")
	}

	defaultQuota := valueObject.NewAccountQuotaWithDefaultValues()
	if createDto.Quota == nil {
		createDto.Quota = &defaultQuota
	}

	accountId, err := accountCmdRepo.Create(createDto)
	if err != nil {
		slog.Error("CreateAccountInfraError", slog.Any("error", err))
		return errors.New("CreateAccountInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateAccount(createDto, accountId)

	return nil
}
