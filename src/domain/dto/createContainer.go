package dto

import (
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type CreateContainer struct {
	AccountId          valueObject.AccountId               `json:"accountId"`
	Hostname           valueObject.Fqdn                    `json:"hostname"`
	ImageAddress       valueObject.ContainerImageAddress   `json:"imageAddress"`
	PortBindings       []valueObject.PortBinding           `json:"portBindings"`
	RestartPolicy      *valueObject.ContainerRestartPolicy `json:"restartPolicy"`
	Entrypoint         *valueObject.ContainerEntrypoint    `json:"entrypoint"`
	ProfileId          *valueObject.ContainerProfileId     `json:"profileId"`
	Envs               []valueObject.ContainerEnv          `json:"envs"`
	AutoCreateMappings bool                                `json:"autoCreateMappings"`
}

func NewCreateContainer(
	accountId valueObject.AccountId,
	hostname valueObject.Fqdn,
	imageAddress valueObject.ContainerImageAddress,
	portBindings []valueObject.PortBinding,
	restartPolicyPtr *valueObject.ContainerRestartPolicy,
	entrypoint *valueObject.ContainerEntrypoint,
	profileId *valueObject.ContainerProfileId,
	envs []valueObject.ContainerEnv,
	autoCreateMappings bool,
) CreateContainer {
	if restartPolicyPtr == nil {
		restartPolicy, _ := valueObject.NewContainerRestartPolicy("unless-stopped")
		restartPolicyPtr = &restartPolicy
	}

	defaultContainerProfileId := entity.DefaultContainerProfile().Id
	if profileId == nil {
		profileId = &defaultContainerProfileId
	}

	return CreateContainer{
		AccountId:          accountId,
		Hostname:           hostname,
		ImageAddress:       imageAddress,
		PortBindings:       portBindings,
		RestartPolicy:      restartPolicyPtr,
		Entrypoint:         entrypoint,
		ProfileId:          profileId,
		Envs:               envs,
		AutoCreateMappings: autoCreateMappings,
	}
}
