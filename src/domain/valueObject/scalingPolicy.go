package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

type ScalingPolicy string

var ValidScalingPolicies = []string{"cpu", "memory", "connection"}

func NewScalingPolicy(value interface{}) (policy ScalingPolicy, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return policy, errors.New("ScalingPolicyMustBeString")
	}

	stringValue = strings.ToLower(stringValue)

	if !slices.Contains(ValidScalingPolicies, stringValue) {
		switch stringValue {
		case "connections", "conn", "conns":
			stringValue = "connection"
		case "mem", "ram":
			stringValue = "memory"
		default:
			return policy, errors.New("InvalidScalingPolicy")
		}
	}

	return ScalingPolicy(stringValue), nil
}

func DefaultScalingPolicy() ScalingPolicy {
	scalingPolicy, _ := NewScalingPolicy("cpu")
	return scalingPolicy
}

func (vo ScalingPolicy) String() string {
	return string(vo)
}
