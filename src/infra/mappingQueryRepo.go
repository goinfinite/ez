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
	dbSvc *db.DatabaseService
}

func NewMappingQueryRepo(dbSvc *db.DatabaseService) *MappingQueryRepo {
	return &MappingQueryRepo{dbSvc: dbSvc}
}

func (repo MappingQueryRepo) Get() ([]entity.Mapping, error) {
	mappingEntities := []entity.Mapping{}

	var mappingModels []dbModel.Mapping

	err := repo.dbSvc.Orm.Model(&dbModel.Mapping{}).
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

func (repo MappingQueryRepo) GetById(id valueObject.MappingId) (entity.Mapping, error) {
	var mapping entity.Mapping

	mappingModel := dbModel.Mapping{ID: uint(id.Get())}

	err := repo.dbSvc.Orm.Model(&mappingModel).
		Preload("Targets").
		First(&mappingModel).Error
	if err != nil {
		return mapping, errors.New("MappingNotFound")
	}

	return mappingModel.ToEntity()
}

func (repo MappingQueryRepo) GetByProtocol(
	protocol valueObject.NetworkProtocol,
) ([]entity.Mapping, error) {
	mappingEntities := []entity.Mapping{}

	var mappingModels []dbModel.Mapping
	err := repo.dbSvc.Orm.Model(dbModel.Mapping{}).
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

func (repo MappingQueryRepo) GetTargetById(
	id valueObject.MappingTargetId,
) (entity.MappingTarget, error) {
	var mappingTarget entity.MappingTarget

	mappingTargetModel := dbModel.MappingTarget{ID: uint(id.Get())}

	err := repo.dbSvc.Orm.Model(&mappingTargetModel).
		First(&mappingTargetModel).Error
	if err != nil {
		return mappingTarget, errors.New("MappingTargetNotFound")
	}

	return mappingTargetModel.ToEntity()
}

func (repo MappingQueryRepo) FindOne(
	hostname *valueObject.Fqdn,
	publicPort valueObject.NetworkPort,
	protocol valueObject.NetworkProtocol,
) (entity.Mapping, error) {
	var mapping entity.Mapping

	mappingModel := dbModel.Mapping{
		PublicPort: uint(publicPort.Get()),
		Protocol:   protocol.String(),
	}

	query := repo.dbSvc.Orm.Model(&mappingModel).Preload("Targets")

	whereHostname := "hostname IS NULL"
	if hostname != nil {
		whereHostname = "hostname = '" + hostname.String() + "'"
	}
	query = query.Where(whereHostname)

	queryResult := query.Find(&mappingModel)
	if queryResult.Error != nil {
		return mapping, errors.New("DbQueryMappingError")
	}

	if queryResult.RowsAffected == 0 {
		return mapping, errors.New("MappingNotFound")
	}

	mappingEntity, err := mappingModel.ToEntity()
	if err != nil {
		return mapping, err
	}

	return mappingEntity, nil
}
