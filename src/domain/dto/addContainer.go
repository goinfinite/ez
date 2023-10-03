package dto

import "github.com/speedianet/sfm/src/domain/valueObject"

type AddContainer struct {
	AccountId         valueObject.AccountId               `json:"accountId"`
	Hostname          valueObject.Fqdn                    `json:"hostname"`
	Image             valueObject.ContainerImgAddress     `json:"image"`
	PortBindings      []valueObject.PortBinding           `json:"portBindings"`
	RestartPolicy     *valueObject.ContainerRestartPolicy `json:"restartPolicy"`
	Entrypoint        *valueObject.ContainerEntrypoint    `json:"entrypoint"`
	ResourceProfileId *valueObject.ResourceProfileId      `json:"resourceProfileId"`
	Envs              []valueObject.ContainerEnv          `json:"envs"`
}

func NewAddContainer(
	accountId valueObject.AccountId,
	hostname valueObject.Fqdn,
	image valueObject.ContainerImgAddress,
	portBindings []valueObject.PortBinding,
	restartPolicy *valueObject.ContainerRestartPolicy,
	entrypoint *valueObject.ContainerEntrypoint,
	resourceProfileId *valueObject.ResourceProfileId,
	envs []valueObject.ContainerEnv,
) AddContainer {
	return AddContainer{
		AccountId:         accountId,
		Hostname:          hostname,
		Image:             image,
		PortBindings:      portBindings,
		RestartPolicy:     restartPolicy,
		Entrypoint:        entrypoint,
		ResourceProfileId: resourceProfileId,
		Envs:              envs,
	}
}
