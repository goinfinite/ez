package infra

import (
	"time"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
	dbModel "github.com/speedianet/control/src/infra/db/model"
)

type ScheduledTaskCmdRepo struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewScheduledTaskCmdRepo(
	persistentDbSvc *db.PersistentDatabaseService,
) *ScheduledTaskCmdRepo {
	return &ScheduledTaskCmdRepo{persistentDbSvc: persistentDbSvc}
}

func (repo *ScheduledTaskCmdRepo) Create(
	createDto dto.CreateScheduledTask,
) error {
	newTaskStatus, _ := valueObject.NewScheduledTaskStatus("pending")

	tags := []string{}
	if len(createDto.Tags) > 0 {
		for _, tag := range createDto.Tags {
			tags = append(tags, tag.String())
		}
	}

	var runAtPtr *time.Time
	if createDto.RunAt != nil {
		runAt := time.Unix(createDto.RunAt.Get(), 0)
		runAtPtr = &runAt
	}

	scheduledTaskModel := dbModel.NewScheduledTask(
		0,
		createDto.Name.String(),
		newTaskStatus.String(),
		createDto.Command.String(),
		tags,
		createDto.TimeoutSecs,
		runAtPtr,
		nil,
		nil,
	)

	return repo.persistentDbSvc.Handler.Create(&scheduledTaskModel).Error
}
