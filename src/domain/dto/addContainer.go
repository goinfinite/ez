package dto

import "github.com/speedianet/sfm/src/domain/valueObject"

type AddContainer struct {
	Hostname      valueObject.Fqdn                    `json:"hostname"`
	Image         valueObject.ContainerImgAddress     `json:"image"`
	PortBindings  *[]valueObject.PortBinding          `json:"portBindings"`
	RestartPolicy *valueObject.ContainerRestartPolicy `json:"restartPolicy"`
	Entrypoint    *valueObject.ContainerEntrypoint    `json:"entrypoint"`
	BaseSpecs     *valueObject.ContainerSpecs         `json:"baseSpecs"`
	MaxSpecs      *valueObject.ContainerSpecs         `json:"maxSpecs"`
	Envs          *[]valueObject.ContainerEnv         `json:"envs"`
}

func NewAddContainer(
	hostname valueObject.Fqdn,
	image valueObject.ContainerImgAddress,
	portBindings *[]valueObject.PortBinding,
	restartPolicy *valueObject.ContainerRestartPolicy,
	entrypoint *valueObject.ContainerEntrypoint,
	baseSpecs *valueObject.ContainerSpecs,
	maxSpecs *valueObject.ContainerSpecs,
	envs *[]valueObject.ContainerEnv,
) AddContainer {
	return AddContainer{
		Hostname:      hostname,
		Image:         image,
		PortBindings:  portBindings,
		RestartPolicy: restartPolicy,
		Entrypoint:    entrypoint,
		BaseSpecs:     baseSpecs,
		MaxSpecs:      maxSpecs,
		Envs:          envs,
	}
}
