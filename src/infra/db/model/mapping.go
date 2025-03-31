package dbModel

import (
	"time"

	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type Mapping struct {
	ID          uint64 `gorm:"primarykey"`
	AccountID   uint64
	AccountName string
	Hostname    *string
	PublicPort  uint
	Protocol    string
	Targets     []MappingTarget
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Mapping) TableName() string {
	return "mappings"
}

func NewMapping(
	id, accountId uint64,
	accountName string,
	hostname *string,
	publicPort uint,
	protocol string,
	targets []MappingTarget,
	createdAt, updatedAt time.Time,
) Mapping {
	mappingModel := Mapping{
		AccountID:   accountId,
		AccountName: accountName,
		Hostname:    hostname,
		PublicPort:  publicPort,
		Protocol:    protocol,
		Targets:     targets,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
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

	accountName, err := valueObject.NewUnixUsername(model.AccountName)
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

	createdAt := valueObject.NewUnixTimeWithGoTime(model.CreatedAt)
	updatedAt := valueObject.NewUnixTimeWithGoTime(model.UpdatedAt)

	return entity.NewMapping(
		mappingId, accountId, accountName, hostnamePtr, port, protocol, targets,
		createdAt, updatedAt,
	), nil
}
