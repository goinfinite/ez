package infra

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/infra/db"
	dbModel "github.com/speedianet/control/src/infra/db/model"
)

type SecurityCmdRepo struct {
	trailDbSvc *db.TrailDatabaseService
}

func NewSecurityCmdRepo(
	trailDbSvc *db.TrailDatabaseService,
) *SecurityCmdRepo {
	return &SecurityCmdRepo{trailDbSvc: trailDbSvc}
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
		accountIdUint := uint(createDto.AccountId.Read())
		accountIdUintPtr = &accountIdUint
	}

	securityEventModel := dbModel.NewSecurityEvent(
		0,
		createDto.Type.String(),
		detailsStrPtr,
		ipAddressStrPtr,
		accountIdUintPtr,
	)

	return repo.trailDbSvc.Handler.Create(&securityEventModel).Error
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
		deleteConditionsMap["account_id"] = deleteDto.AccountId.Read()
	}

	return repo.trailDbSvc.Handler.
		Where(deleteConditionsMap).
		Delete(&dbModel.SecurityEvent{}).
		Error
}
