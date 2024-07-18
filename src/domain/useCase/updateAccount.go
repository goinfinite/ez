package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func UpdateAccount(
	accountQueryRepo repository.AccountQueryRepo,
	accountCmdRepo repository.AccountCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	updateDto dto.UpdateAccount,
) error {
	_, err := accountQueryRepo.ReadById(updateDto.AccountId)
	if err != nil {
		return errors.New("AccountNotFound")
	}

	if updateDto.Password != nil {
		err = accountCmdRepo.UpdatePassword(updateDto.AccountId, *updateDto.Password)
		if err != nil {
			log.Printf("UpdateAccountPasswordError: %s", err)
			return errors.New("UpdateAccountPasswordInfraError")
		}

		recordCode, _ := valueObject.NewActivityRecordCode("AccountPasswordUpdated")
		CreateSecurityActivityRecord(
			activityRecordCmdRepo, &recordCode, &updateDto.IpAddress,
			&updateDto.OperatorAccountId, &updateDto.AccountId, nil,
		)
	}

	if updateDto.Quota == nil {
		return nil
	}

	err = accountCmdRepo.UpdateQuota(updateDto.AccountId, *updateDto.Quota)
	if err != nil {
		log.Printf("UpdateAccountQuotaError: %s", err)
		return errors.New("UpdateAccountQuotaInfraError")
	}

	recordCode, _ := valueObject.NewActivityRecordCode("AccountQuotaUpdated")
	CreateSecurityActivityRecord(
		activityRecordCmdRepo, &recordCode, &updateDto.IpAddress,
		&updateDto.OperatorAccountId, &updateDto.AccountId, nil,
	)

	return nil
}
