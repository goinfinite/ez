package dbModel

import (
	"time"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type Mapping struct {
	ID         uint `gorm:"primarykey"`
	AccountID  uint
	Hostname   *string
	PublicPort uint
	Protocol   string
	Targets    []MappingTarget
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (Mapping) TableName() string {
	return "mappings"
}

func NewMapping(
	id uint,
	accountId uint,
	hostname *string,
	publicPort uint,
	protocol string,
	targets []MappingTarget,
	createdAt time.Time,
	updatedAt time.Time,
) Mapping {
	mappingModel := Mapping{
		AccountID:  accountId,
		Hostname:   hostname,
		PublicPort: publicPort,
		Protocol:   protocol,
		Targets:    targets,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}

	if id != 0 {
		mappingModel.ID = id
	}

	return mappingModel
}

func (model Mapping) ToEntity() (entity.Mapping, error) {
	var mapping entity.Mapping

	mappingId, err := valueObject.NewMappingId(model.ID)
	if err != nil {
		return mapping, err
	}

	accountId, err := valueObject.NewAccountId(model.AccountID)
	if err != nil {
		return mapping, err
	}

	var hostnamePtr *valueObject.Fqdn
	if model.Hostname != nil {
		hostname, err := valueObject.NewFqdn(*model.Hostname)
		if err != nil {
			return mapping, err
		}
		hostnamePtr = &hostname
	}

	port, err := valueObject.NewNetworkPort(model.PublicPort)
	if err != nil {
		return mapping, err
	}

	protocol, err := valueObject.NewNetworkProtocol(model.Protocol)
	if err != nil {
		return mapping, err
	}

	var targets []entity.MappingTarget
	for _, target := range model.Targets {
		targetEntity, err := target.ToEntity()
		if err != nil {
			return mapping, err
		}
		targets = append(targets, targetEntity)
	}

	createdAt := valueObject.UnixTime(model.CreatedAt.Unix())
	updatedAt := valueObject.UnixTime(model.UpdatedAt.Unix())

	return entity.NewMapping(
		mappingId,
		accountId,
		hostnamePtr,
		port,
		protocol,
		targets,
		createdAt,
		updatedAt,
	), nil
}

func (Mapping) AddDtoToModel(addDto dto.AddMapping) Mapping {
	var hostnamePtr *string
	if addDto.Hostname != nil {
		hostnameStr := addDto.Hostname.String()
		hostnamePtr = &hostnameStr
	}

	return NewMapping(
		0,
		uint(addDto.AccountId.Get()),
		hostnamePtr,
		uint(addDto.PublicPort.Get()),
		addDto.Protocol.String(),
		[]MappingTarget{},
		time.Now(),
		time.Now(),
	)
}
