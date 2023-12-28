package dbModel

import (
	"time"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type Mapping struct {
	ID        uint `gorm:"primarykey"`
	AccountID uint
	Hostname  *string
	Port      uint
	Protocol  string
	Targets   []MappingTarget
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Mapping) TableName() string {
	return "mappings"
}

func NewMapping(
	id uint,
	accountId uint,
	hostname *string,
	port uint,
	protocol string,
	targets []MappingTarget,
) Mapping {
	mappingStruct := Mapping{
		AccountID: accountId,
		Hostname:  hostname,
		Port:      port,
		Protocol:  protocol,
		Targets:   targets,
	}

	if id != 0 {
		mappingStruct.ID = id
	}

	return mappingStruct
}

func NewMappingWithAddDto(addDto dto.AddMapping) Mapping {
	var hostnamePtr *string
	if addDto.Hostname != nil {
		hostnameStr := addDto.Hostname.String()
		hostnamePtr = &hostnameStr
	}

	mappingStruct := NewMapping(
		0,
		uint(addDto.AccountId.Get()),
		hostnamePtr,
		uint(addDto.Port.Get()),
		addDto.Protocol.String(),
		[]MappingTarget{},
	)

	return mappingStruct
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

	port, err := valueObject.NewNetworkPort(model.Port)
	if err != nil {
		return mapping, err
	}

	protocol, err := valueObject.NewNetworkProtocol(model.Protocol)
	if err != nil {
		return mapping, err
	}

	var targets []valueObject.MappingTarget
	for _, target := range model.Targets {
		targetVo, err := target.ToValueObject()
		if err != nil {
			return mapping, err
		}
		targets = append(targets, targetVo)
	}

	return entity.NewMapping(
		mappingId,
		accountId,
		hostnamePtr,
		port,
		protocol,
		targets,
	), nil
}
