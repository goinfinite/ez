package dto

import "github.com/speedianet/control/src/domain/valueObject"

type UpdateContainerProfile struct {
	AccountId              valueObject.AccountId             `json:"accountId"`
	ProfileId              valueObject.ContainerProfileId    `json:"profileId"`
	Name                   *valueObject.ContainerProfileName `json:"name"`
	BaseSpecs              *valueObject.ContainerSpecs       `json:"baseSpecs"`
	MaxSpecs               *valueObject.ContainerSpecs       `json:"maxSpecs"`
	ScalingPolicy          *valueObject.ScalingPolicy        `json:"scalingPolicy"`
	ScalingThreshold       *uint                             `json:"scalingThreshold"`
	ScalingMaxDurationSecs *uint                             `json:"scalingMaxDurationSecs"`
	ScalingIntervalSecs    *uint                             `json:"scalingIntervalSecs"`
	HostMinCapacityPercent *valueObject.HostMinCapacity      `json:"hostMinCapacityPercent"`
	OperatorAccountId      valueObject.AccountId             `json:"-"`
	OperatorIpAddress      valueObject.IpAddress             `json:"-"`
}

func NewUpdateContainerProfile(
	accountId valueObject.AccountId,
	profileId valueObject.ContainerProfileId,
	name *valueObject.ContainerProfileName,
	baseSpecs *valueObject.ContainerSpecs,
	maxSpecs *valueObject.ContainerSpecs,
	scalingPolicy *valueObject.ScalingPolicy,
	scalingThreshold, scalingMaxDurationSecs, scalingIntervalSecs *uint,
	hostMinCapacityPercent *valueObject.HostMinCapacity,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) UpdateContainerProfile {
	return UpdateContainerProfile{
		AccountId:              accountId,
		ProfileId:              profileId,
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
