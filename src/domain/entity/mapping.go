package entity

import "github.com/speedianet/control/src/domain/valueObject"

type Mapping struct {
	Id         valueObject.MappingId       `json:"id"`
	AccountId  valueObject.AccountId       `json:"accountId"`
	Hostname   *valueObject.Fqdn           `json:"hostname"`
	PublicPort valueObject.NetworkPort     `json:"publicPort"`
	Protocol   valueObject.NetworkProtocol `json:"protocol"`
	Targets    []MappingTarget             `json:"targets"`
}

func NewMapping(
	id valueObject.MappingId,
	accountId valueObject.AccountId,
	hostname *valueObject.Fqdn,
	publicPort valueObject.NetworkPort,
	protocol valueObject.NetworkProtocol,
	targets []MappingTarget,
) Mapping {
	return Mapping{
		Id:         id,
		AccountId:  accountId,
		Hostname:   hostname,
		PublicPort: publicPort,
		Protocol:   protocol,
		Targets:    targets,
	}
}
