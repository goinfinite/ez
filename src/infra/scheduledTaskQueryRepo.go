package infra

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/infra/db"
	dbHelper "github.com/goinfinite/ez/src/infra/db/helper"
	dbModel "github.com/goinfinite/ez/src/infra/db/model"
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
	if readDto.TaskName != nil {
		scheduledTaskModel.Name = readDto.TaskName.String()
	}
	if readDto.TaskStatus != nil {
		scheduledTaskModel.Status = readDto.TaskStatus.String()
	}

	dbQuery := repo.persistentDbSvc.Handler.Model(scheduledTaskModel).Where(&scheduledTaskModel)
	if len(readDto.TaskTags) == 0 {
		dbQuery = dbQuery.Preload("Tags")
	} else {
		tagsStrSlice := []string{}
		for _, taskTag := range readDto.TaskTags {
			tagsStrSlice = append(tagsStrSlice, taskTag.String())
		}
		dbQuery = dbQuery.
			Joins("JOIN scheduled_tasks_tags ON scheduled_tasks_tags.scheduled_task_id = scheduled_tasks.id").
			Where("scheduled_tasks_tags.tag IN (?)", tagsStrSlice)
	}
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

	paginatedDbQuery, responsePagination, err := dbHelper.PaginationQueryBuilder(
		dbQuery, readDto.Pagination,
	)
	if err != nil {
		return responseDto, errors.New("PaginationQueryBuilderError: " + err.Error())
	}

	scheduledTaskModels := []dbModel.ScheduledTask{}
	err = paginatedDbQuery.Find(&scheduledTaskModels).Error
	if err != nil {
		return responseDto, errors.New("FindScheduledTasksError: " + err.Error())
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

	return dto.ReadScheduledTasksResponse{
		Pagination: responsePagination,
		Tasks:      scheduledTaskEntities,
	}, nil
}
