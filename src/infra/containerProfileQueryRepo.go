package infra

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
	"github.com/goinfinite/ez/src/infra/db"
	dbModel "github.com/goinfinite/ez/src/infra/db/model"
)

type ContainerProfileQueryRepo struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewContainerProfileQueryRepo(
	persistentDbSvc *db.PersistentDatabaseService,
) *ContainerProfileQueryRepo {
	return &ContainerProfileQueryRepo{persistentDbSvc: persistentDbSvc}
}

func (repo *ContainerProfileQueryRepo) Read() ([]entity.ContainerProfile, error) {
	profileEntities := []entity.ContainerProfile{}

	var profileModels []dbModel.ContainerProfile
	err := repo.persistentDbSvc.Handler.
		Model(&dbModel.ContainerProfile{}).
		Find(&profileModels).Error
	if err != nil {
		return profileEntities, errors.New("QueryContainerProfilesError: " + err.Error())
	}

	for _, profileModel := range profileModels {
		profileEntity, err := profileModel.ToEntity()
		if err != nil {
			slog.Debug(
				"ModelToEntityError",
				slog.Uint64("id", profileModel.ID),
				slog.Any("error", err),
			)
			continue
		}

		profileEntities = append(profileEntities, profileEntity)
	}

	return profileEntities, nil
}

func (repo *ContainerProfileQueryRepo) ReadById(
	profileId valueObject.ContainerProfileId,
) (profileEntity entity.ContainerProfile, err error) {
	profileModel := dbModel.ContainerProfile{ID: profileId.Uint64()}
	err = repo.persistentDbSvc.Handler.
		Model(&profileModel).
		First(&profileModel).Error
	if err != nil {
		return profileEntity, err
	}

	return profileModel.ToEntity()
}
