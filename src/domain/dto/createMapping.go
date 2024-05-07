package dto

import "github.com/speedianet/control/src/domain/valueObject"

type CreateMapping struct {
	AccountId    valueObject.AccountId            `json:"accountId"`
	Hostname     *valueObject.Fqdn                `json:"hostname"`
	PublicPort   valueObject.NetworkPort          `json:"publicPort"`
	Protocol     valueObject.NetworkProtocol      `json:"protocol"`
	SourcePath   *valueObject.MappingPath         `json:"sourcePath"`
	MatchPattern *valueObject.MappingMatchPattern `json:"matchPattern"`
	TargetPath   *valueObject.MappingPath         `json:"targetPath"`
	ContainerIds []valueObject.ContainerId        `json:"containerIds"`
}

func NewCreateMapping(
	accountId valueObject.AccountId,
	hostname *valueObject.Fqdn,
	publicPort valueObject.NetworkPort,
	protocol valueObject.NetworkProtocol,
	sourcePath *valueObject.MappingPath,
	matchPattern *valueObject.MappingMatchPattern,
	targetPath *valueObject.MappingPath,
	containerIds []valueObject.ContainerId,
) CreateMapping {
	return CreateMapping{
		AccountId:    accountId,
		Hostname:     hostname,
		PublicPort:   publicPort,
		Protocol:     protocol,
		SourcePath:   sourcePath,
		MatchPattern: matchPattern,
		TargetPath:   targetPath,
		ContainerIds: containerIds,
	}
}
