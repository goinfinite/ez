package valueObject

import "errors"

type ContainerRestartPolicy string

const (
	no            ContainerRestartPolicy = "no"
	onFailure     ContainerRestartPolicy = "on-failure"
	always        ContainerRestartPolicy = "always"
	unlessStopped ContainerRestartPolicy = "unless-stopped"
)

func NewContainerRestartPolicy(value string) (ContainerRestartPolicy, error) {
	crp := ContainerRestartPolicy(value)
	if !crp.isValid() {
		return "", errors.New("InvalidContainerRestartPolicy")
	}
	return crp, nil
}

func NewContainerRestartPolicyPanic(value string) ContainerRestartPolicy {
	crp := ContainerRestartPolicy(value)
	if !crp.isValid() {
		panic("InvalidContainerRestartPolicy")
	}
	return crp
}

func (crp ContainerRestartPolicy) isValid() bool {
	switch crp {
	case no, onFailure, always, unlessStopped:
		return true
	default:
		return false
	}
}

func (crp ContainerRestartPolicy) String() string {
	return string(crp)
}
