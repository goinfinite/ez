package dbModel

import (
	"time"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type Mapping struct {
	ID           uint `gorm:"primarykey"`
	AccountID    uint
	Hostname     *string
	PublicPort   uint
	Protocol     string
	Path         *string
	MatchPattern *string
	Targets      []MappingTarget
	CreatedAt    time.Time
	UpdatedAt    time.Time
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
	path *string,
	matchPattern *string,
	targets []MappingTarget,
	createdAt time.Time,
	updatedAt time.Time,
) Mapping {
	mappingModel := Mapping{
		AccountID:    accountId,
		Hostname:     hostname,
		PublicPort:   publicPort,
		Protocol:     protocol,
		Path:         path,
		MatchPattern: matchPattern,
		Targets:      targets,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
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

	var pathPtr *valueObject.MappingPath
	if model.Path != nil {
		path, err := valueObject.NewMappingPath(*model.Path)
		if err != nil {
			return mapping, err
		}
		pathPtr = &path
	}

	var matchPatternPtr *valueObject.MappingMatchPattern
	if model.MatchPattern != nil {
		matchPattern, err := valueObject.NewMappingMatchPattern(*model.MatchPattern)
		if err != nil {
			return mapping, err
		}
		matchPatternPtr = &matchPattern
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
		pathPtr,
		matchPatternPtr,
		targets,
		createdAt,
		updatedAt,
	), nil
}

func (Mapping) CreateDtoToModel(createDto dto.CreateMapping) Mapping {
	var hostnamePtr *string
	if createDto.Hostname != nil {
		hostnameStr := createDto.Hostname.String()
		hostnamePtr = &hostnameStr
	}

	var pathPtr *string
	if createDto.Path != nil {
		pathStr := createDto.Path.String()
		pathPtr = &pathStr
	}

	var matchPatternPtr *string
	if createDto.MatchPattern != nil {
		matchPatternStr := createDto.MatchPattern.String()
		matchPatternPtr = &matchPatternStr
	}

	return NewMapping(
		0,
		uint(createDto.AccountId.Get()),
		hostnamePtr,
		uint(createDto.PublicPort.Get()),
		createDto.Protocol.String(),
		pathPtr,
		matchPatternPtr,
		[]MappingTarget{},
		time.Now(),
		time.Now(),
	)
}
