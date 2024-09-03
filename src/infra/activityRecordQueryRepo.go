package infra

import (
	"log/slog"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/infra/db"
	dbModel "github.com/speedianet/control/src/infra/db/model"
)

type ActivityRecordQueryRepo struct {
	trailDbSvc *db.TrailDatabaseService
}

func NewActivityRecordQueryRepo(
	trailDbSvc *db.TrailDatabaseService,
) *ActivityRecordQueryRepo {
	return &ActivityRecordQueryRepo{trailDbSvc: trailDbSvc}
}

func (repo *ActivityRecordQueryRepo) Read(
	readDto dto.ReadActivityRecords,
) ([]entity.ActivityRecord, error) {
	activityRecordEvents := []entity.ActivityRecord{}

	readModel := dbModel.ActivityRecord{}
	if readDto.RecordLevel != nil {
		recordLevelStr := readDto.RecordLevel.String()
		readModel.RecordLevel = recordLevelStr
	}

	if readDto.RecordCode != nil {
		recordCodeStr := readDto.RecordCode.String()
		readModel.RecordCode = recordCodeStr
	}

	if readDto.OperatorAccountId != nil {
		operatorAccountId := readDto.OperatorAccountId.Uint64()
		readModel.OperatorAccountId = &operatorAccountId
	}

	if readDto.OperatorIpAddress != nil {
		operatorIpAddressStr := readDto.OperatorIpAddress.String()
		readModel.OperatorIpAddress = &operatorIpAddressStr
	}

	if readDto.AccountId != nil {
		accountId := readDto.AccountId.Uint64()
		readModel.AccountId = &accountId
	}

	if readDto.ContainerId != nil {
		containerIdStr := readDto.ContainerId.String()
		readModel.ContainerId = &containerIdStr
	}

	if readDto.ContainerProfileId != nil {
		containerProfileId := readDto.ContainerProfileId.Uint64()
		readModel.ContainerProfileId = &containerProfileId
	}

	if readDto.ContainerImageId != nil {
		containerImageId := readDto.ContainerImageId.String()
		readModel.ContainerImageId = &containerImageId
	}

	if readDto.MappingId != nil {
		mappingId := readDto.MappingId.Uint64()
		readModel.MappingId = &mappingId
	}

	if readDto.ScheduledTaskId != nil {
		scheduledTaskId := readDto.ScheduledTaskId.Uint64()
		readModel.ScheduledTaskId = &scheduledTaskId
	}

	dbQuery := repo.trailDbSvc.Handler.Where(&readModel)
	if readDto.CreatedBeforeAt != nil {
		dbQuery = dbQuery.Where("created_at < ?", readDto.CreatedBeforeAt.GetAsGoTime())
	}
	if readDto.CreatedAfterAt != nil {
		dbQuery = dbQuery.Where("created_at > ?", readDto.CreatedAfterAt.GetAsGoTime())
	}

	activityRecordEventModels := []dbModel.ActivityRecord{}
	err := dbQuery.Find(&activityRecordEventModels).Error
	if err != nil {
		return activityRecordEvents, err
	}

	for _, activityRecordEventModel := range activityRecordEventModels {
		activityRecordEvent, err := activityRecordEventModel.ToEntity()
		if err != nil {
			slog.Debug(
				"ModelToEntityError",
				slog.Uint64("id", activityRecordEventModel.ID),
				slog.Any("error", err),
			)
			continue
		}
		activityRecordEvents = append(activityRecordEvents, activityRecordEvent)
	}

	return activityRecordEvents, nil
}
