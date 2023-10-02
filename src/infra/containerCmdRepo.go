package infra

import (
	"time"

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
		"--hostname",
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

	baseSpecsParams := []string{
		"--cpus",
		addContainer.BaseSpecs.CpuCores.String(),
		"--memory",
		addContainer.BaseSpecs.MemoryBytes.String(),
	}
	runParams = append(runParams, baseSpecsParams...)

	if addContainer.MaxSpecs != nil {
		maxSpecsParams := []string{
			"--annotation",
			"speedia/max-cpu=" + addContainer.MaxSpecs.CpuCores.String(),
			"--annotation",
			"speedia/max-memory=" + addContainer.MaxSpecs.MemoryBytes.String(),
		}
		runParams = append(runParams, maxSpecsParams...)
	}

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
	time.Sleep(1 * time.Second)

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
	currentContainer, err := ContainerQueryRepo{}.GetById(
		updateContainer.AccountId,
		updateContainer.ContainerId,
	)
	if err != nil {
		return err
	}

	shouldUpdateStatus := updateContainer.Status != currentContainer.Status
	if shouldUpdateStatus {
		actionToPerform := "start"
		if !updateContainer.Status {
			actionToPerform = "stop"
		}

		_, err := infraHelper.RunCmdAsUser(
			updateContainer.AccountId,
			"podman",
			actionToPerform,
			updateContainer.ContainerId.String(),
		)
		if err != nil {
			return err
		}
	}

	if updateContainer.BaseSpecs != nil {
		_, err := infraHelper.RunCmdAsUser(
			updateContainer.AccountId,
			"podman",
			"update",
			"--cpus",
			updateContainer.BaseSpecs.CpuCores.String(),
			"--memory",
			updateContainer.BaseSpecs.MemoryBytes.String(),
			updateContainer.ContainerId.String(),
		)
		if err != nil {
			return err
		}
	}

	// TODO: Update podman max specs annotation

	return nil
}

func (repo ContainerCmdRepo) Delete(
	accId valueObject.AccountId,
	containerId valueObject.ContainerId,
) error {
	_, err := infraHelper.RunCmdAsUser(
		accId,
		"podman",
		"rm",
		"--force",
		containerId.String(),
	)
	if err != nil {
		return err
	}

	return nil
}
