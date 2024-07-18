package useCase

import (
	"errors"
	"log"

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
		log.Printf("CreateAccountError: %s", err)
		return errors.New("CreateAccountInfraError")
	}

	recordCode, _ := valueObject.NewActivityRecordCode("AccountCreated")
	CreateSecurityActivityRecord(
		activityRecordCmdRepo, &recordCode, &createDto.IpAddress,
		&createDto.OperatorAccountId, &accountId, &createDto.Username,
	)

	return nil
}
