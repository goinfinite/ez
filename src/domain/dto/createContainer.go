package dto

import (
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type CreateContainer struct {
	AccountId            valueObject.AccountId               `json:"accountId"`
	Hostname             valueObject.Fqdn                    `json:"hostname"`
	ImageAddress         valueObject.ContainerImageAddress   `json:"imageAddress"`
	ImageId              *valueObject.ContainerImageId       `json:"imageId,omitempty"`
	PortBindings         []valueObject.PortBinding           `json:"portBindings"`
	RestartPolicy        *valueObject.ContainerRestartPolicy `json:"restartPolicy,omitempty"`
	Entrypoint           *valueObject.ContainerEntrypoint    `json:"entrypoint,omitempty"`
	ProfileId            *valueObject.ContainerProfileId     `json:"profileId,omitempty"`
	Envs                 []valueObject.ContainerEnv          `json:"envs"`
	LaunchScript         *valueObject.LaunchScript           `json:"launchScript,omitempty"`
	AutoCreateMappings   bool                                `json:"autoCreateMappings"`
	UseImageExposedPorts bool                                `json:"useImageExposedPorts"`
	ExistingContainerId  *valueObject.ContainerId            `json:"existingContainerId,omitempty"`
	OperatorAccountId    valueObject.AccountId               `json:"-"`
	OperatorIpAddress    valueObject.IpAddress               `json:"-"`
}

func NewCreateContainer(
	accountId valueObject.AccountId,
	hostname valueObject.Fqdn,
	imageAddress valueObject.ContainerImageAddress,
	imageId *valueObject.ContainerImageId,
	portBindings []valueObject.PortBinding,
	restartPolicyPtr *valueObject.ContainerRestartPolicy,
	entrypoint *valueObject.ContainerEntrypoint,
	profileId *valueObject.ContainerProfileId,
	envs []valueObject.ContainerEnv,
	launchScript *valueObject.LaunchScript,
	autoCreateMappings bool,
	useImageExposedPorts bool,
	existingContainerId *valueObject.ContainerId,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
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
		AccountId:            accountId,
		Hostname:             hostname,
		ImageAddress:         imageAddress,
		ImageId:              imageId,
		PortBindings:         portBindings,
		RestartPolicy:        restartPolicyPtr,
		Entrypoint:           entrypoint,
		ProfileId:            profileId,
		Envs:                 envs,
		LaunchScript:         launchScript,
		AutoCreateMappings:   autoCreateMappings,
		UseImageExposedPorts: useImageExposedPorts,
		ExistingContainerId:  existingContainerId,
		OperatorAccountId:    operatorAccountId,
		OperatorIpAddress:    operatorIpAddress,
	}
}
