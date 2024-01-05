package dto

import "github.com/speedianet/control/src/domain/valueObject"

type AddMapping struct {
	AccountId  valueObject.AccountId              `json:"accountId"`
	Hostname   *valueObject.Fqdn                  `json:"hostname"`
	PublicPort valueObject.NetworkPort            `json:"publicPort"`
	Protocol   valueObject.NetworkProtocol        `json:"protocol"`
	Targets    []AddMappingTargetWithoutMappingId `json:"targets"`
}

func NewAddMapping(
	accountId valueObject.AccountId,
	hostname *valueObject.Fqdn,
	publicPort valueObject.NetworkPort,
	protocol valueObject.NetworkProtocol,
	targets []AddMappingTargetWithoutMappingId,
) AddMapping {
	return AddMapping{
		AccountId:  accountId,
		Hostname:   hostname,
		PublicPort: publicPort,
		Protocol:   protocol,
		Targets:    targets,
	}
}
