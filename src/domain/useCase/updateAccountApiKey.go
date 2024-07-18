package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func UpdateAccountApiKey(
	accountQueryRepo repository.AccountQueryRepo,
	accountCmdRepo repository.AccountCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	updateDto dto.UpdateAccount,
) (newKey valueObject.AccessTokenValue, err error) {
	_, err = accountQueryRepo.ReadById(updateDto.AccountId)
	if err != nil {
		return newKey, errors.New("AccountNotFound")
	}

	newKey, err = accountCmdRepo.UpdateApiKey(updateDto.AccountId)
	if err != nil {
		log.Printf("UpdateAccountApiKeyError: %s", err)
		return newKey, errors.New("UpdateAccountApiKeyInfraError")
	}

	recordCode, _ := valueObject.NewActivityRecordCode("AccountApiKeyUpdated")
	CreateSecurityActivityRecord(
		activityRecordCmdRepo, &recordCode, &updateDto.IpAddress,
		&updateDto.OperatorAccountId, &updateDto.AccountId, nil,
	)

	return newKey, nil
}
