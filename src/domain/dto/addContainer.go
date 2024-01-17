package dto

import (
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type AddContainer struct {
	AccountId          valueObject.AccountId               `json:"accountId"`
	Hostname           valueObject.Fqdn                    `json:"hostname"`
	ImageAddr          valueObject.ContainerImgAddress     `json:"imageAddr"`
	PortBindings       []valueObject.PortBinding           `json:"portBindings"`
	RestartPolicy      *valueObject.ContainerRestartPolicy `json:"restartPolicy"`
	Entrypoint         *valueObject.ContainerEntrypoint    `json:"entrypoint"`
	ProfileId          *valueObject.ContainerProfileId     `json:"profileId"`
	Envs               []valueObject.ContainerEnv          `json:"envs"`
	AutoCreateMappings bool                                `json:"autoCreateMappings"`
}

func NewAddContainer(
	accountId valueObject.AccountId,
	hostname valueObject.Fqdn,
	imageAddr valueObject.ContainerImgAddress,
	portBindings []valueObject.PortBinding,
	restartPolicyPtr *valueObject.ContainerRestartPolicy,
	entrypoint *valueObject.ContainerEntrypoint,
	profileId *valueObject.ContainerProfileId,
	envs []valueObject.ContainerEnv,
	autoCreateMappings bool,
) AddContainer {
	if restartPolicyPtr == nil {
		restartPolicy, _ := valueObject.NewContainerRestartPolicy("unless-stopped")
		restartPolicyPtr = &restartPolicy
	}

	defaultContainerProfileId := entity.DefaultContainerProfile().Id
	if profileId == nil {
		profileId = &defaultContainerProfileId
	}

	return AddContainer{
		AccountId:          accountId,
		Hostname:           hostname,
		ImageAddr:          imageAddr,
		PortBindings:       portBindings,
		RestartPolicy:      restartPolicyPtr,
		Entrypoint:         entrypoint,
		ProfileId:          profileId,
		Envs:               envs,
		AutoCreateMappings: autoCreateMappings,
	}
}
