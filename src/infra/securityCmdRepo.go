package infra

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/infra/db"
	dbModel "github.com/speedianet/control/src/infra/db/model"
)

type SecurityCmdRepo struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewSecurityCmdRepo(
	persistentDbSvc *db.PersistentDatabaseService,
) *SecurityCmdRepo {
	return &SecurityCmdRepo{persistentDbSvc: persistentDbSvc}
}

func (repo *SecurityCmdRepo) CreateEvent(createDto dto.CreateSecurityEvent) error {
	var detailsStrPtr *string
	if createDto.Details != nil {
		detailsStr := createDto.Details.String()
		detailsStrPtr = &detailsStr
	}

	var ipAddressStrPtr *string
	if createDto.IpAddress != nil {
		ipAddressStr := createDto.IpAddress.String()
		ipAddressStrPtr = &ipAddressStr
	}

	var accountIdUintPtr *uint
	if createDto.AccountId != nil {
		accountIdUint := uint(createDto.AccountId.Get())
		accountIdUintPtr = &accountIdUint
	}

	securityEventModel := dbModel.NewSecurityEvent(
		0,
		createDto.Type.String(),
		detailsStrPtr,
		ipAddressStrPtr,
		accountIdUintPtr,
	)

	return repo.persistentDbSvc.Handler.Create(&securityEventModel).Error
}

func (repo *SecurityCmdRepo) DeleteEvents(deleteDto dto.DeleteSecurityEvents) error {
	deleteConditionsMap := map[string]interface{}{}

	if deleteDto.Type != nil {
		deleteConditionsMap["type"] = deleteDto.Type.String()
	}

	if deleteDto.IpAddress != nil {
		deleteConditionsMap["ip_address"] = deleteDto.IpAddress.String()
	}

	if deleteDto.AccountId != nil {
		deleteConditionsMap["account_id"] = deleteDto.AccountId.Get()
	}

	return repo.persistentDbSvc.Handler.
		Where(deleteConditionsMap).
		Delete(&dbModel.SecurityEvent{}).
		Error
}
