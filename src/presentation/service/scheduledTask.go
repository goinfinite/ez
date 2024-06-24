package service

import (
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
)

type ScheduledTaskService struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewScheduledTaskService(
	persistentDbSvc *db.PersistentDatabaseService,
) *ScheduledTaskService {
	return &ScheduledTaskService{
		persistentDbSvc: persistentDbSvc,
	}
}

func (service *ScheduledTaskService) Read() ServiceOutput {
	scheduledTaskQueryRepo := infra.NewScheduledTaskQueryRepo(service.persistentDbSvc)
	scheduledTasksList, err := useCase.GetScheduledTasks(scheduledTaskQueryRepo)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, scheduledTasksList)
}
