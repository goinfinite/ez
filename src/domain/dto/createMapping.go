package dto

import "github.com/speedianet/control/src/domain/valueObject"

type CreateMapping struct {
	AccountId    valueObject.AccountId            `json:"accountId"`
	Hostname     *valueObject.Fqdn                `json:"hostname"`
	PublicPort   valueObject.NetworkPort          `json:"publicPort"`
	Protocol     valueObject.NetworkProtocol      `json:"protocol"`
	Path         *valueObject.MappingPath         `json:"path"`
	MatchPattern *valueObject.MappingMatchPattern `json:"matchPattern"`
	ContainerIds []valueObject.ContainerId        `json:"containerIds"`
}

func NewCreateMapping(
	accountId valueObject.AccountId,
	hostname *valueObject.Fqdn,
	publicPort valueObject.NetworkPort,
	protocol valueObject.NetworkProtocol,
	path *valueObject.MappingPath,
	matchPattern *valueObject.MappingMatchPattern,
	containerIds []valueObject.ContainerId,
) CreateMapping {
	return CreateMapping{
		AccountId:    accountId,
		Hostname:     hostname,
		PublicPort:   publicPort,
		Protocol:     protocol,
		Path:         path,
		MatchPattern: matchPattern,
		ContainerIds: containerIds,
	}
}
