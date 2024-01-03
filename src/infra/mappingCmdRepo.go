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
	return repo.dbSvc.Orm.Create(&mappingModel).Error
}

func (repo MappingCmdRepo) AddTarget(addDto dto.AddMappingTarget) error {
	model := dbModel.NewMappingTarget(
		0,
		uint(addDto.MappingId),
		addDto.ContainerId.String(),
		nil,
		nil,
	)

	if addDto.Port != nil {
		portUint := uint(addDto.Port.Get())
		model.Port = &portUint
	}

	if addDto.Protocol != nil {
		protocolStr := addDto.Protocol.String()
		model.Protocol = &protocolStr
	}

	createResult := repo.dbSvc.Orm.Create(&model)
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
