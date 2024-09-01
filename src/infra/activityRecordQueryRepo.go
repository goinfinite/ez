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
	if readDto.Level != nil {
		levelStr := readDto.Level.String()
		readModel.Level = levelStr
	}

	if readDto.Code != nil {
		codeStr := readDto.Code.String()
		readModel.Code = &codeStr
	}

	if readDto.Message != nil {
		messageStr := readDto.Message.String()
		readModel.Message = &messageStr
	}

	if readDto.IpAddress != nil {
		ipAddressStr := readDto.IpAddress.String()
		readModel.IpAddress = &ipAddressStr
	}

	if readDto.OperatorAccountId != nil {
		operatorAccountId := readDto.OperatorAccountId.Uint64()
		readModel.OperatorAccountId = &operatorAccountId
	}

	if readDto.TargetAccountId != nil {
		targetAccountId := readDto.TargetAccountId.Uint64()
		readModel.TargetAccountId = &targetAccountId
	}

	if readDto.Username != nil {
		usernameStr := readDto.Username.String()
		readModel.Username = &usernameStr
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

	dbQuery := repo.trailDbSvc.Handler.Where(&readModel)
	if readDto.CreatedAt != nil {
		dbQuery = dbQuery.Where("created_at >= ?", readDto.CreatedAt.GetAsGoTime())
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
