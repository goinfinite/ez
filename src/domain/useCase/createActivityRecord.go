package useCase

import (
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/valueObject"
)

func CreateSecurityActivityRecord(
	code *valueObject.ActivityRecordCode,
	ipAddress *valueObject.IpAddress,
	operatorAccountId *valueObject.AccountId,
	targetAccountId *valueObject.AccountId,
	username *valueObject.Username,
) {
	recordLevel, _ := valueObject.NewActivityRecordLevel("SEC")

	createDto := dto.CreateActivityRecord{
		Level:             recordLevel,
		Code:              code,
		OperatorAccountId: operatorAccountId,
		TargetAccountId:   targetAccountId,
		IpAddress:         ipAddress,
		Username:          username,
	}

	log.Print(createDto)
}
