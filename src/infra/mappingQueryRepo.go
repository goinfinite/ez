package infra

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
	dbModel "github.com/speedianet/control/src/infra/db/model"
)

type MappingQueryRepo struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewMappingQueryRepo(
	persistentDbSvc *db.PersistentDatabaseService,
) *MappingQueryRepo {
	return &MappingQueryRepo{persistentDbSvc: persistentDbSvc}
}

func (repo *MappingQueryRepo) Get() ([]entity.Mapping, error) {
	mappingEntities := []entity.Mapping{}

	var mappingModels []dbModel.Mapping

	err := repo.persistentDbSvc.Handler.Model(&dbModel.Mapping{}).
		Preload("Targets").
		Find(&mappingModels).Error
	if err != nil {
		return mappingEntities, errors.New("DbQueryMappingError")
	}

	for _, mappingModel := range mappingModels {
		mappingEntity, err := mappingModel.ToEntity()
		if err != nil {
			log.Printf("MappingModelToEntityError: %v", err.Error())
			continue
		}

		mappingEntities = append(mappingEntities, mappingEntity)
	}

	return mappingEntities, nil
}

func (repo *MappingQueryRepo) GetById(id valueObject.MappingId) (entity.Mapping, error) {
	var mapping entity.Mapping

	mappingModel := dbModel.Mapping{ID: uint(id.Get())}

	err := repo.persistentDbSvc.Handler.Model(&mappingModel).
		Preload("Targets").
		First(&mappingModel).Error
	if err != nil {
		return mapping, errors.New("MappingNotFound")
	}

	return mappingModel.ToEntity()
}

func (repo *MappingQueryRepo) GetByProtocol(
	protocol valueObject.NetworkProtocol,
) ([]entity.Mapping, error) {
	mappingEntities := []entity.Mapping{}

	var mappingModels []dbModel.Mapping
	err := repo.persistentDbSvc.Handler.Model(dbModel.Mapping{}).
		Preload("Targets").
		Where("protocol = ?", protocol.String()).
		Find(&mappingModels).Error
	if err != nil {
		return mappingEntities, errors.New("GetMappingsFromDatabaseError")
	}

	for _, mappingModel := range mappingModels {
		mappingEntity, err := mappingModel.ToEntity()
		if err != nil {
			log.Printf("MappingModelToEntityError: %v", err.Error())
			continue
		}

		mappingEntities = append(mappingEntities, mappingEntity)
	}

	return mappingEntities, nil
}

func (repo *MappingQueryRepo) GetTargetById(
	id valueObject.MappingTargetId,
) (entity.MappingTarget, error) {
	var mappingTarget entity.MappingTarget

	mappingTargetModel := dbModel.MappingTarget{ID: uint(id.Get())}

	err := repo.persistentDbSvc.Handler.Model(&mappingTargetModel).
		First(&mappingTargetModel).Error
	if err != nil {
		return mappingTarget, errors.New("MappingTargetNotFound")
	}

	return mappingTargetModel.ToEntity()
}

func (repo *MappingQueryRepo) GetTargetsByContainerId(
	containerId valueObject.ContainerId,
) ([]entity.MappingTarget, error) {
	mappingTargets := []entity.MappingTarget{}

	var mappingTargetModels []dbModel.MappingTarget

	err := repo.persistentDbSvc.Handler.
		Model(&dbModel.MappingTarget{}).
		Where("container_id = ?", containerId.String()).
		Find(&mappingTargetModels).Error
	if err != nil {
		return mappingTargets, errors.New("GetTargetsFromDatabaseError")
	}

	for _, mappingTargetModel := range mappingTargetModels {
		mappingTargetEntity, err := mappingTargetModel.ToEntity()
		if err != nil {
			log.Printf("MappingTargetModelToEntityError: %v", err.Error())
			continue
		}

		mappingTargets = append(mappingTargets, mappingTargetEntity)
	}

	return mappingTargets, nil
}
