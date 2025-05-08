package infra

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
	"github.com/goinfinite/ez/src/infra/db"
	dbModel "github.com/goinfinite/ez/src/infra/db/model"
)

type MappingQueryRepo struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewMappingQueryRepo(
	persistentDbSvc *db.PersistentDatabaseService,
) *MappingQueryRepo {
	return &MappingQueryRepo{persistentDbSvc: persistentDbSvc}
}

func (repo *MappingQueryRepo) Read() ([]entity.Mapping, error) {
	mappingEntities := []entity.Mapping{}

	var mappingModels []dbModel.Mapping
	err := repo.persistentDbSvc.Handler.
		Preload("Account").
		Preload("Targets").
		Find(&mappingModels).Error
	if err != nil {
		return mappingEntities, errors.New("ReadMappingsFromDatabaseError: " + err.Error())
	}

	for _, mappingModel := range mappingModels {
		mappingEntity, err := mappingModel.ToEntity()
		if err != nil {
			slog.Debug(
				"MappingModelToEntityError",
				slog.Uint64("mappingId", mappingModel.ID),
				slog.String("error", err.Error()),
			)
			continue
		}

		mappingEntities = append(mappingEntities, mappingEntity)
	}

	return mappingEntities, nil
}

func (repo *MappingQueryRepo) ReadByContainerId(
	containerId valueObject.ContainerId,
) ([]entity.Mapping, error) {
	mappingEntities := []entity.Mapping{}

	var mappingModels []dbModel.Mapping
	err := repo.persistentDbSvc.Handler.
		Preload("Account").
		Preload("Targets", "container_id = ?", containerId.String()).
		Model(&dbModel.Mapping{}).
		Find(&mappingModels).Error
	if err != nil {
		return mappingEntities, errors.New("ReadMappingsFromDatabaseError: " + err.Error())
	}

	for _, mappingModel := range mappingModels {
		if len(mappingModel.Targets) == 0 {
			continue
		}

		mappingEntity, err := mappingModel.ToEntity()
		if err != nil {
			slog.Debug(
				"MappingModelToEntityError",
				slog.Uint64("mappingId", mappingModel.ID),
				slog.String("error", err.Error()),
			)
			continue
		}

		mappingEntities = append(mappingEntities, mappingEntity)
	}

	return mappingEntities, nil
}

func (repo *MappingQueryRepo) ReadById(
	id valueObject.MappingId,
) (mappingEntity entity.Mapping, err error) {
	var mappingModel dbModel.Mapping
	err = repo.persistentDbSvc.Handler.
		Preload("Account").
		Preload("Targets").
		Where("id = ?", id.Uint64()).
		First(&mappingModel).Error
	if err != nil {
		return mappingEntity, errors.New("MappingNotFound")
	}

	return mappingModel.ToEntity()
}

func (repo *MappingQueryRepo) GetByProtocol(
	protocol valueObject.NetworkProtocol,
) ([]entity.Mapping, error) {
	mappingEntities := []entity.Mapping{}

	var mappingModels []dbModel.Mapping
	err := repo.persistentDbSvc.Handler.
		Preload("Account").
		Preload("Targets").
		Where("protocol = ?", protocol.String()).
		Find(&mappingModels).Error
	if err != nil {
		return mappingEntities, errors.New("ReadMappingsFromDatabaseError: " + err.Error())
	}

	for _, mappingModel := range mappingModels {
		mappingEntity, err := mappingModel.ToEntity()
		if err != nil {
			slog.Debug(
				"MappingModelToEntityError",
				slog.Uint64("mappingId", mappingModel.ID),
				slog.String("error", err.Error()),
			)
			continue
		}

		mappingEntities = append(mappingEntities, mappingEntity)
	}

	return mappingEntities, nil
}

func (repo *MappingQueryRepo) ReadTargetById(
	id valueObject.MappingTargetId,
) (mappingTargetEntity entity.MappingTarget, err error) {
	var mappingTargetModel dbModel.MappingTarget

	err = repo.persistentDbSvc.Handler.
		Where("id = ?", id.Uint64()).
		First(&mappingTargetModel).Error
	if err != nil {
		return mappingTargetEntity, errors.New("MappingTargetNotFound")
	}

	return mappingTargetModel.ToEntity()
}

func (repo *MappingQueryRepo) ReadTargetsByContainerId(
	containerId valueObject.ContainerId,
) ([]entity.MappingTarget, error) {
	mappingTargets := []entity.MappingTarget{}

	var mappingTargetModels []dbModel.MappingTarget

	err := repo.persistentDbSvc.Handler.
		Where("container_id = ?", containerId.String()).
		Find(&mappingTargetModels).Error
	if err != nil {
		return mappingTargets, errors.New("ReadTargetsFromDatabaseError: " + err.Error())
	}

	for _, mappingTargetModel := range mappingTargetModels {
		mappingTargetEntity, err := mappingTargetModel.ToEntity()
		if err != nil {
			slog.Debug(
				"MappingTargetModelToEntityError",
				slog.Uint64("mappingTargetId", mappingTargetModel.ID),
				slog.String("error", err.Error()),
			)
			continue
		}

		mappingTargets = append(mappingTargets, mappingTargetEntity)
	}

	return mappingTargets, nil
}
