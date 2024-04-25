package valueObject

import (
	"errors"
	"strings"

	"golang.org/x/exp/slices"
)

type ContainerRestartPolicy string

var ValidContainerRestartPolicies = []string{
	"no",
	"on-failure",
	"always",
	"unless-stopped",
}

func NewContainerRestartPolicy(value interface{}) (ContainerRestartPolicy, error) {
	stringValue, assertOk := value.(string)
	if !assertOk {
		return "", errors.New("ContainerRestartPolicyMustBeString")
	}

	stringValue = strings.TrimSpace(stringValue)
	stringValue = strings.ToLower(stringValue)

	if !slices.Contains(ValidContainerRestartPolicies, stringValue) {
		return "", errors.New("UnknownContainerRestartPolicy")
	}

	return ContainerRestartPolicy(stringValue), nil
}

func NewContainerRestartPolicyPanic(value string) ContainerRestartPolicy {
	crp, err := NewContainerRestartPolicy(value)
	if err != nil {
		panic(err)
	}
	return crp
}

func (crp ContainerRestartPolicy) String() string {
	return string(crp)
}
