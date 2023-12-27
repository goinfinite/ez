package dto

import "github.com/speedianet/control/src/domain/valueObject"

type AddMapping struct {
	Hostname       *valueObject.Fqdn           `json:"hostname"`
	Port           valueObject.NetworkPort     `json:"port"`
	Protocol       valueObject.NetworkProtocol `json:"protocol"`
	MappingTargets []valueObject.MappingTarget `json:"mappingTargets"`
}

func NewAddMapping(
	hostname *valueObject.Fqdn,
	port valueObject.NetworkPort,
	protocol valueObject.NetworkProtocol,
	mappingTargets []valueObject.MappingTarget,
) AddMapping {
	return AddMapping{
		Hostname:       hostname,
		Port:           port,
		Protocol:       protocol,
		MappingTargets: mappingTargets,
	}
}
