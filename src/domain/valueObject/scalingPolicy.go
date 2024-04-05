package valueObject

import (
	"errors"
	"strings"

	"golang.org/x/exp/slices"
)

type ScalingPolicy string

var ValidScalingPolicies = []string{
	"connection",
	"cpu",
	"memory",
}

func NewScalingPolicy(value string) (ScalingPolicy, error) {
	value = strings.TrimSpace(value)
	value = strings.ToLower(value)

	if !slices.Contains(ValidScalingPolicies, value) {
		return "", errors.New("InvalidScalingPolicy")
	}
	return ScalingPolicy(value), nil
}

func DefaultScalingPolicy() ScalingPolicy {
	scalingPolicy, _ := NewScalingPolicy("cpu")
	return scalingPolicy
}

func NewScalingPolicyPanic(value string) ScalingPolicy {
	scalingPolicy, err := NewScalingPolicy(value)
	if err != nil {
		panic(err)
	}
	return scalingPolicy
}

func (sp ScalingPolicy) String() string {
	return string(sp)
}
