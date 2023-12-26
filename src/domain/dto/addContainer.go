package dto

import "github.com/speedianet/control/src/domain/valueObject"

type AddContainer struct {
	AccountId     valueObject.AccountId               `json:"accountId"`
	Hostname      valueObject.Fqdn                    `json:"hostname"`
	ImgAddr       valueObject.ContainerImgAddress     `json:"imgAddr"`
	PortBindings  []valueObject.PortBinding           `json:"portBindings"`
	RestartPolicy *valueObject.ContainerRestartPolicy `json:"restartPolicy"`
	Entrypoint    *valueObject.ContainerEntrypoint    `json:"entrypoint"`
	ProfileId     *valueObject.ContainerProfileId     `json:"profileId"`
	Envs          []valueObject.ContainerEnv          `json:"envs"`
}

func NewAddContainer(
	accountId valueObject.AccountId,
	hostname valueObject.Fqdn,
	imgAddr valueObject.ContainerImgAddress,
	portBindings []valueObject.PortBinding,
	restartPolicy *valueObject.ContainerRestartPolicy,
	entrypoint *valueObject.ContainerEntrypoint,
	profileId *valueObject.ContainerProfileId,
	envs []valueObject.ContainerEnv,
) AddContainer {
	return AddContainer{
		AccountId:     accountId,
		Hostname:      hostname,
		ImgAddr:       imgAddr,
		PortBindings:  portBindings,
		RestartPolicy: restartPolicy,
		Entrypoint:    entrypoint,
		ProfileId:     profileId,
		Envs:          envs,
	}
}
