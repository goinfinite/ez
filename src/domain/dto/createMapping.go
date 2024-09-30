package dto

import "github.com/goinfinite/ez/src/domain/valueObject"

type CreateMapping struct {
	AccountId         valueObject.AccountId       `json:"accountId"`
	Hostname          *valueObject.Fqdn           `json:"hostname"`
	PublicPort        valueObject.NetworkPort     `json:"publicPort"`
	Protocol          valueObject.NetworkProtocol `json:"protocol"`
	ContainerIds      []valueObject.ContainerId   `json:"containerIds"`
	OperatorAccountId valueObject.AccountId       `json:"-"`
	OperatorIpAddress valueObject.IpAddress       `json:"-"`
}

func NewCreateMapping(
	accountId valueObject.AccountId,
	hostname *valueObject.Fqdn,
	publicPort valueObject.NetworkPort,
	protocol valueObject.NetworkProtocol,
	containerIds []valueObject.ContainerId,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) CreateMapping {
	return CreateMapping{
		AccountId:         accountId,
		Hostname:          hostname,
		PublicPort:        publicPort,
		Protocol:          protocol,
		ContainerIds:      containerIds,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
