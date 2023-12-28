package entity

import "github.com/speedianet/control/src/domain/valueObject"

type Mapping struct {
	Id             valueObject.MappingId       `json:"id"`
	AccountId      valueObject.AccountId       `json:"accountId"`
	Hostname       *valueObject.Fqdn           `json:"hostname"`
	Port           valueObject.NetworkPort     `json:"port"`
	Protocol       valueObject.NetworkProtocol `json:"protocol"`
	MappingTargets []valueObject.MappingTarget `json:"mappingTargets"`
}

func NewMapping(
	id valueObject.MappingId,
	accountId valueObject.AccountId,
	hostname *valueObject.Fqdn,
	port valueObject.NetworkPort,
	protocol valueObject.NetworkProtocol,
	mappingTargets []valueObject.MappingTarget,
) Mapping {
	return Mapping{
		Id:             id,
		AccountId:      accountId,
		Hostname:       hostname,
		Port:           port,
		Protocol:       protocol,
		MappingTargets: mappingTargets,
	}
}
