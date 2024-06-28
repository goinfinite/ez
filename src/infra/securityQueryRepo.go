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

	getConditionsMap := map[string]interface{}{}

	if readDto.Type != nil {
		getConditionsMap["type"] = readDto.Type.String()
	}

	if readDto.IpAddress != nil {
		getConditionsMap["ip_address"] = readDto.IpAddress.String()
	}

	if readDto.AccountId != nil {
		getConditionsMap["account_id"] = readDto.AccountId.Get()
	}

	securityEventModels := []dbModel.SecurityEvent{}
	err := repo.persistentDbSvc.Handler.
		Where(getConditionsMap).
		Find(&securityEventModels).
		Error
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
