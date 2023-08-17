package entity

import "github.com/speedianet/sfm/src/domain/valueObject"

type Container struct {
	Id               valueObject.ContainerId            `json:"id"`
	Hostname         valueObject.Fqdn                   `json:"hostname"`
	Status           bool                               `json:"status"`
	Image            valueObject.ContainerImgAddress    `json:"image"`
	ImageHash        valueObject.Hash                   `json:"imageHash"`
	PrivateIpAddress valueObject.IpAddress              `json:"privateIp"`
	PortBindings     []valueObject.PortBinding          `json:"portBindings"`
	RestartPolicy    valueObject.ContainerRestartPolicy `json:"restartPolicy"`
	Entrypoint       valueObject.ContainerEntrypoint    `json:"entrypoint"`
	CreatedAt        valueObject.UnixTime               `json:"createdAt"`
	StartedAt        valueObject.UnixTime               `json:"startedAt"`
	BaseSpecs        valueObject.ContainerSpecs         `json:"baseSpecs"`
	MaxSpecs         valueObject.ContainerSpecs         `json:"maxSpecs"`
	Envs             []valueObject.ContainerEnv         `json:"envs"`
}

func NewContainer(
	id valueObject.ContainerId,
	hostname valueObject.Fqdn,
	status bool,
	image valueObject.ContainerImgAddress,
	imageHash valueObject.Hash,
	privateIpAddress valueObject.IpAddress,
	portBindings []valueObject.PortBinding,
	restartPolicy valueObject.ContainerRestartPolicy,
	entrypoint valueObject.ContainerEntrypoint,
	createdAt valueObject.UnixTime,
	startedAt valueObject.UnixTime,
	baseSpecs valueObject.ContainerSpecs,
	maxSpecs valueObject.ContainerSpecs,
	envs []valueObject.ContainerEnv,
) Container {
	return Container{
		Id:               id,
		Hostname:         hostname,
		Status:           status,
		Image:            image,
		ImageHash:        imageHash,
		PrivateIpAddress: privateIpAddress,
		PortBindings:     portBindings,
		RestartPolicy:    restartPolicy,
		Entrypoint:       entrypoint,
		CreatedAt:        createdAt,
		StartedAt:        startedAt,
		BaseSpecs:        baseSpecs,
		MaxSpecs:         maxSpecs,
		Envs:             envs,
	}
}
