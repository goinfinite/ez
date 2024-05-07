package entity

import (
	"github.com/speedianet/control/src/domain/valueObject"
)

type Container struct {
	Id            valueObject.ContainerId            `json:"id"`
	AccountId     valueObject.AccountId              `json:"accountId"`
	Hostname      valueObject.Fqdn                   `json:"hostname"`
	Status        bool                               `json:"status"`
	ImageAddress  valueObject.ContainerImageAddress  `json:"imageAddress"`
	ImageHash     valueObject.Hash                   `json:"imageHash"`
	PortBindings  []valueObject.PortBinding          `json:"portBindings"`
	RestartPolicy valueObject.ContainerRestartPolicy `json:"restartPolicy"`
	RestartCount  uint64                             `json:"restartCount"`
	Entrypoint    *valueObject.ContainerEntrypoint   `json:"entrypoint"`
	ProfileId     valueObject.ContainerProfileId     `json:"profileId"`
	Envs          []valueObject.ContainerEnv         `json:"envs"`
	CreatedAt     valueObject.UnixTime               `json:"createdAt"`
	UpdatedAt     valueObject.UnixTime               `json:"updatedAt"`
	StartedAt     *valueObject.UnixTime              `json:"startedAt"`
	StoppedAt     *valueObject.UnixTime              `json:"stoppedAt"`
}

func NewContainer(
	id valueObject.ContainerId,
	accountId valueObject.AccountId,
	hostname valueObject.Fqdn,
	status bool,
	imageAddress valueObject.ContainerImageAddress,
	imageHash valueObject.Hash,
	portBindings []valueObject.PortBinding,
	restartPolicy valueObject.ContainerRestartPolicy,
	restartCount uint64,
	entrypoint *valueObject.ContainerEntrypoint,
	profileId valueObject.ContainerProfileId,
	envs []valueObject.ContainerEnv,
	createdAt valueObject.UnixTime,
	updatedAt valueObject.UnixTime,
	startedAt *valueObject.UnixTime,
	stoppedAt *valueObject.UnixTime,
) Container {
	return Container{
		Id:            id,
		AccountId:     accountId,
		Hostname:      hostname,
		Status:        status,
		ImageAddress:  imageAddress,
		ImageHash:     imageHash,
		PortBindings:  portBindings,
		RestartPolicy: restartPolicy,
		RestartCount:  restartCount,
		Entrypoint:    entrypoint,
		ProfileId:     profileId,
		Envs:          envs,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
		StartedAt:     startedAt,
		StoppedAt:     stoppedAt,
	}
}

func (container *Container) IsRunning() bool {
	return container.Status
}
