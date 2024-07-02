package infra

import (
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/infra/db"
	dbModel "github.com/speedianet/control/src/infra/db/model"
)

type SecurityQueryRepo struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewSecurityQueryRepo(
	persistentDbSvc *db.PersistentDatabaseService,
) *SecurityQueryRepo {
	return &SecurityQueryRepo{persistentDbSvc: persistentDbSvc}
}

func (repo *SecurityQueryRepo) ReadEvents(
	readDto dto.ReadSecurityEvents,
) ([]entity.SecurityEvent, error) {
	securityEvents := []entity.SecurityEvent{}

	dbQuery := repo.persistentDbSvc.Handler.Model(&dbModel.SecurityEvent{})
	if readDto.Type != nil {
		dbQuery = dbQuery.Where("type = ?", readDto.Type.String())
	}

	if readDto.IpAddress != nil {
		dbQuery = dbQuery.Where("ip_address = ?", readDto.IpAddress.String())
	}

	if readDto.AccountId != nil {
		dbQuery = dbQuery.Where("account_id = ?", readDto.AccountId.Read())
	}

	if readDto.CreatedAt != nil {
		dbQuery = dbQuery.Where("created_at >= ?", readDto.CreatedAt.GetAsGoTime())
	}

	securityEventModels := []dbModel.SecurityEvent{}
	err := dbQuery.Find(&securityEventModels).Error
	if err != nil {
		return securityEvents, err
	}

	for _, securityEventModel := range securityEventModels {
		securityEvent, err := securityEventModel.ToEntity()
		if err != nil {
			log.Printf("EventModelToEntityError: %v", err.Error())
			continue
		}
		securityEvents = append(securityEvents, securityEvent)
	}

	return securityEvents, nil
}
