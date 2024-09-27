package infra

import (
	"log/slog"

	"github.com/speedianet/control/src/domain/dto"
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

func (repo *ScheduledTaskQueryRepo) Read(
	readDto dto.ReadScheduledTasksRequest,
) (responseDto dto.ReadScheduledTasksResponse, err error) {
	scheduledTaskEntities := []entity.ScheduledTask{}

	scheduledTaskModel := dbModel.ScheduledTask{}
	if readDto.TaskId != nil {
		scheduledTaskModel.ID = readDto.TaskId.Uint64()
	}

	if readDto.TaskStatus != nil {
		scheduledTaskModel.Status = readDto.TaskStatus.String()
	}

	taskTagsModels := []dbModel.ScheduledTaskTag{}
	if readDto.TaskTags != nil {
		for _, taskTag := range readDto.TaskTags {
			taskTagModel := dbModel.ScheduledTaskTag{
				Tag: taskTag.String(),
			}
			taskTagsModels = append(taskTagsModels, taskTagModel)
		}
		scheduledTaskModel.Tags = taskTagsModels
	}

	dbQuery := repo.persistentDbSvc.Handler.Where(&scheduledTaskModel)
	if readDto.StartedBeforeAt != nil {
		dbQuery = dbQuery.Where("started_at < ?", readDto.StartedBeforeAt.GetAsGoTime())
	}
	if readDto.StartedAfterAt != nil {
		dbQuery = dbQuery.Where("started_at > ?", readDto.StartedAfterAt.GetAsGoTime())
	}
	if readDto.FinishedBeforeAt != nil {
		dbQuery = dbQuery.Where("finished_at < ?", readDto.FinishedBeforeAt.GetAsGoTime())
	}
	if readDto.FinishedAfterAt != nil {
		dbQuery = dbQuery.Where("finished_at > ?", readDto.FinishedAfterAt.GetAsGoTime())
	}
	if readDto.CreatedBeforeAt != nil {
		dbQuery = dbQuery.Where("created_at < ?", readDto.CreatedBeforeAt.GetAsGoTime())
	}
	if readDto.CreatedAfterAt != nil {
		dbQuery = dbQuery.Where("created_at > ?", readDto.CreatedAfterAt.GetAsGoTime())
	}

	scheduledTaskModels := []dbModel.ScheduledTask{}
	err = repo.persistentDbSvc.Handler.
		Preload("Tags").
		Find(&scheduledTaskModels).Error
	if err != nil {
		return responseDto, err
	}

	for _, scheduledTaskModel := range scheduledTaskModels {
		scheduledTaskEntity, err := scheduledTaskModel.ToEntity()
		if err != nil {
			slog.Debug(
				"ModelToEntityError",
				slog.Uint64("id", scheduledTaskModel.ID),
				slog.Any("error", err),
			)
			continue
		}
		scheduledTaskEntities = append(scheduledTaskEntities, scheduledTaskEntity)
	}

	responsePagination := dto.Pagination{
		PageNumber:   readDto.Pagination.PageNumber,
		ItemsPerPage: readDto.Pagination.ItemsPerPage,
	}

	return dto.ReadScheduledTasksResponse{
		Pagination: responsePagination,
		Tasks:      scheduledTaskEntities,
	}, nil
}
