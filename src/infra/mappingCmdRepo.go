package infra

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
	dbModel "github.com/speedianet/control/src/infra/db/model"
)

type MappingCmdRepo struct {
	dbSvc *db.DatabaseService
}

func NewMappingCmdRepo(dbSvc *db.DatabaseService) *MappingCmdRepo {
	return &MappingCmdRepo{dbSvc: dbSvc}
}

func (repo MappingCmdRepo) Add(mappingDto dto.AddMapping) error {
	mappingModel := dbModel.NewMappingWithAddDto(mappingDto)

	createResult := repo.dbSvc.Orm.Create(&mappingModel)
	if createResult.Error != nil {
		return createResult.Error
	}

	mappingId, err := valueObject.NewMappingId(mappingModel.ID)
	if err != nil {
		return err
	}

	return repo.AddTargets(mappingId, mappingDto.Targets)
}

func (repo MappingCmdRepo) AddTargets(
	mappingId valueObject.MappingId,
	targets []valueObject.MappingTarget,
) error {
	targetModels := []dbModel.MappingTarget{}
	for _, targetVo := range targets {
		model := dbModel.NewMappingTarget(
			0,
			uint(mappingId.Get()),
			targetVo.ContainerId.String(),
			nil,
			nil,
		)

		if targetVo.Port != nil {
			portUint := uint(targetVo.Port.Get())
			model.Port = &portUint
		}

		if targetVo.Protocol != nil {
			protocolStr := targetVo.Protocol.String()
			model.Protocol = &protocolStr
		}

		targetModels = append(targetModels, model)
	}

	createResult := repo.dbSvc.Orm.Create(&targetModels)
	if createResult.Error != nil {
		return createResult.Error
	}

	return nil
}

func (repo MappingCmdRepo) Delete(id valueObject.MappingId) error {
	ormSvc := repo.dbSvc.Orm

	err := ormSvc.Delete(dbModel.MappingTarget{}, "mapping_id = ?", id.Get()).Error
	if err != nil {
		return err
	}

	return ormSvc.Delete(dbModel.Mapping{}, id.Get()).Error
}
