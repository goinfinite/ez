package valueObject

import (
	"errors"
	"strings"

	"slices"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

type ContainerRestartPolicy string

var ValidContainerRestartPolicies = []string{
	"no",
	"never",
	"on-failure",
	"always",
	"unless-stopped",
}

func NewContainerRestartPolicy(value interface{}) (ContainerRestartPolicy, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("ContainerRestartPolicyMustBeString")
	}

	stringValue = strings.TrimSpace(stringValue)
	stringValue = strings.ToLower(stringValue)

	if !slices.Contains(ValidContainerRestartPolicies, stringValue) {
		return "", errors.New("UnknownContainerRestartPolicy")
	}

	switch stringValue {
	case "never":
		stringValue = "no"
	case "unless-stopped":
		stringValue = "always"
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
