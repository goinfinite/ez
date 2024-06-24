package infra

import (
	"log"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/infra/db"
	dbModel "github.com/speedianet/control/src/infra/db/model"
)

type ScheduledTaskQueryRepo struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewScheduledTaskQueryRepo(
	persistentDbSvc *db.PersistentDatabaseService,
) *ScheduledTaskQueryRepo {
	return &ScheduledTaskQueryRepo{persistentDbSvc: persistentDbSvc}
}

func (repo *ScheduledTaskQueryRepo) Get() ([]entity.ScheduledTask, error) {
	scheduledTasks := []entity.ScheduledTask{}

	scheduledTaskModels := []dbModel.ScheduledTask{}
	err := repo.persistentDbSvc.Handler.
		Find(&scheduledTaskModels).Error
	if err != nil {
		return scheduledTasks, err
	}

	for _, scheduledTaskModel := range scheduledTaskModels {
		scheduledTaskEntity, err := scheduledTaskModel.ToEntity()
		if err != nil {
			log.Printf("[%s] %s", scheduledTaskModel.ID, err.Error())
			continue
		}
		scheduledTasks = append(scheduledTasks, scheduledTaskEntity)
	}

	return scheduledTasks, nil
}
