package infra

import (
	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/valueObject"
	infraHelper "github.com/speedianet/sfm/src/infra/helper"
)

type ContainerCmdRepo struct {
}

func (repo ContainerCmdRepo) Add(addContainer dto.AddContainer) error {
	runParams := []string{
		"run",
		"--detach",
		"--name",
		addContainer.Hostname.String(),
		"--env",
		"VIRTUAL_HOST=" + addContainer.Hostname.String(),
	}

	if addContainer.Envs != nil {
		envFlags := []string{}
		for _, env := range *addContainer.Envs {
			envFlags = append(envFlags, "--env")
			envFlags = append(envFlags, env.String())
		}

		runParams = append(runParams, envFlags...)
	}

	specsParams := []string{
		"--cpus",
		addContainer.BaseSpecs.CpuCores.String(),
		"--memory",
		addContainer.BaseSpecs.MemoryBytes.String(),
	}
	runParams = append(runParams, specsParams...)

	if addContainer.RestartPolicy != nil {
		runParams = append(runParams, "--restart", addContainer.RestartPolicy.String())
	}

	if addContainer.Entrypoint != nil {
		runParams = append(runParams, "--entrypoint", addContainer.Entrypoint.String())
	}

	if addContainer.PortBindings != nil {
		portBindingsParams := []string{}
		for _, portBinding := range *addContainer.PortBindings {
			portBindingsParams = append(portBindingsParams, "--publish")
			portBindingsString := portBinding.HostPort.String() +
				":" + portBinding.ContainerPort.String()
			if portBinding.GetProtocol().String() != "" {
				portBindingsString += "/" + portBinding.GetProtocol().String()
			}
			portBindingsParams = append(portBindingsParams, portBindingsString)
		}

		runParams = append(runParams, portBindingsParams...)
	}

	runParams = append(runParams, addContainer.Image.String())

	err := infraHelper.EnableLingering(addContainer.AccountId)
	if err != nil {
		return err
	}

	_, err = infraHelper.RunCmdAsUser(
		addContainer.AccountId,
		"podman",
		runParams...,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repo ContainerCmdRepo) Update(
	updateContainer dto.UpdateContainer,
) error {
	return nil
}

func (repo ContainerCmdRepo) Delete(
	accId valueObject.AccountId,
	containerId valueObject.ContainerId,
) error {
	return nil
}
