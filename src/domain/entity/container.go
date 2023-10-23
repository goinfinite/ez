package entity

import "github.com/speedianet/sfm/src/domain/valueObject"

type Container struct {
	Id               valueObject.ContainerId            `json:"id"`
	AccountId        valueObject.AccountId              `json:"accountId"`
	Hostname         valueObject.Fqdn                   `json:"hostname"`
	Status           bool                               `json:"status"`
	Image            valueObject.ContainerImgAddress    `json:"image"`
	ImageHash        valueObject.Hash                   `json:"imageHash"`
	PrivateIpAddress valueObject.IpAddress              `json:"privateIp"`
	PortBindings     []valueObject.PortBinding          `json:"portBindings"`
	RestartPolicy    valueObject.ContainerRestartPolicy `json:"restartPolicy"`
	Entrypoint       valueObject.ContainerEntrypoint    `json:"entrypoint"`
	CreatedAt        valueObject.UnixTime               `json:"createdAt"`
	StartedAt        *valueObject.UnixTime              `json:"startedAt"`
	ProfileId        valueObject.ContainerProfileId     `json:"profileId"`
	Envs             []valueObject.ContainerEnv         `json:"envs"`
}

func NewContainer(
	id valueObject.ContainerId,
	accountId valueObject.AccountId,
	hostname valueObject.Fqdn,
	status bool,
	image valueObject.ContainerImgAddress,
	imageHash valueObject.Hash,
	privateIpAddress valueObject.IpAddress,
	portBindings []valueObject.PortBinding,
	restartPolicy valueObject.ContainerRestartPolicy,
	entrypoint valueObject.ContainerEntrypoint,
	createdAt valueObject.UnixTime,
	startedAt *valueObject.UnixTime,
	profileId valueObject.ContainerProfileId,
	envs []valueObject.ContainerEnv,
) Container {
	return Container{
		Id:               id,
		AccountId:        accountId,
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
		ProfileId:        profileId,
		Envs:             envs,
	}
}
