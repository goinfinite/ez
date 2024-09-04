package dto

import "github.com/speedianet/control/src/domain/valueObject"

type CreateContainerProfile struct {
	AccountId              valueObject.AccountId            `json:"accountId"`
	Name                   valueObject.ContainerProfileName `json:"name"`
	BaseSpecs              valueObject.ContainerSpecs       `json:"baseSpecs"`
	MaxSpecs               *valueObject.ContainerSpecs      `json:"maxSpecs,omitempty"`
	ScalingPolicy          *valueObject.ScalingPolicy       `json:"scalingPolicy,omitempty"`
	ScalingThreshold       *uint                            `json:"scalingThreshold,omitempty"`
	ScalingMaxDurationSecs *uint                            `json:"scalingMaxDurationSecs,omitempty"`
	ScalingIntervalSecs    *uint                            `json:"scalingIntervalSecs,omitempty"`
	HostMinCapacityPercent *valueObject.HostMinCapacity     `json:"hostMinCapacityPercent,omitempty"`
	OperatorAccountId      valueObject.AccountId            `json:"-"`
	OperatorIpAddress      valueObject.IpAddress            `json:"-"`
}

func NewCreateContainerProfile(
	accountId valueObject.AccountId,
	name valueObject.ContainerProfileName,
	baseSpecs valueObject.ContainerSpecs,
	maxSpecs *valueObject.ContainerSpecs,
	scalingPolicy *valueObject.ScalingPolicy,
	scalingThreshold, scalingMaxDurationSecs, scalingIntervalSecs *uint,
	hostMinCapacityPercent *valueObject.HostMinCapacity,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) CreateContainerProfile {
	return CreateContainerProfile{
		AccountId:              accountId,
		Name:                   name,
		BaseSpecs:              baseSpecs,
		MaxSpecs:               maxSpecs,
		ScalingPolicy:          scalingPolicy,
		ScalingThreshold:       scalingThreshold,
		ScalingMaxDurationSecs: scalingMaxDurationSecs,
		ScalingIntervalSecs:    scalingIntervalSecs,
		HostMinCapacityPercent: hostMinCapacityPercent,
		OperatorAccountId:      operatorAccountId,
		OperatorIpAddress:      operatorIpAddress,
	}
}
