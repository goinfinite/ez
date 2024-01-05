package entity

import "github.com/speedianet/control/src/domain/valueObject"

type Container struct {
	Id            valueObject.ContainerId            `json:"id"`
	AccountId     valueObject.AccountId              `json:"accountId"`
	Hostname      valueObject.Fqdn                   `json:"hostname"`
	Status        bool                               `json:"status"`
	ImageAddr     valueObject.ContainerImgAddress    `json:"imageAddr"`
	ImageHash     valueObject.Hash                   `json:"imageHash"`
	PortBindings  []valueObject.PortBinding          `json:"portBindings"`
	RestartPolicy valueObject.ContainerRestartPolicy `json:"restartPolicy"`
	RestartCount  uint64                             `json:"restartCount"`
	Entrypoint    *valueObject.ContainerEntrypoint   `json:"entrypoint"`
	CreatedAt     valueObject.UnixTime               `json:"createdAt"`
	StartedAt     *valueObject.UnixTime              `json:"startedAt"`
	ProfileId     valueObject.ContainerProfileId     `json:"profileId"`
	Envs          []valueObject.ContainerEnv         `json:"envs"`
}

func NewContainer(
	id valueObject.ContainerId,
	accountId valueObject.AccountId,
	hostname valueObject.Fqdn,
	status bool,
	imageAddr valueObject.ContainerImgAddress,
	imageHash valueObject.Hash,
	portBindings []valueObject.PortBinding,
	restartPolicy valueObject.ContainerRestartPolicy,
	restartCount uint64,
	entrypoint *valueObject.ContainerEntrypoint,
	createdAt valueObject.UnixTime,
	startedAt *valueObject.UnixTime,
	profileId valueObject.ContainerProfileId,
	envs []valueObject.ContainerEnv,
) Container {
	return Container{
		Id:            id,
		AccountId:     accountId,
		Hostname:      hostname,
		Status:        status,
		ImageAddr:     imageAddr,
		ImageHash:     imageHash,
		PortBindings:  portBindings,
		RestartPolicy: restartPolicy,
		RestartCount:  restartCount,
		Entrypoint:    entrypoint,
		CreatedAt:     createdAt,
		StartedAt:     startedAt,
		ProfileId:     profileId,
		Envs:          envs,
	}
}
