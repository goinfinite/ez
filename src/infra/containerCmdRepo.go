package infra

import (
	"time"

	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/entity"
	"github.com/speedianet/sfm/src/domain/valueObject"
	infraHelper "github.com/speedianet/sfm/src/infra/helper"
)

type ContainerCmdRepo struct {
}

func (repo ContainerCmdRepo) getBaseSpecs(
	resourceProfileId valueObject.ResourceProfileId,
) (valueObject.ContainerSpecs, error) {
	resourceProfile, err := ResourceProfileQueryRepo{}.GetById(
		resourceProfileId,
	)
	if err != nil {
		return valueObject.ContainerSpecs{}, err
	}

	return resourceProfile.BaseSpecs, nil
}

func (repo ContainerCmdRepo) Add(addContainer dto.AddContainer) error {
	containerName := addContainer.ResourceProfileId.String() +
		"-" + addContainer.Hostname.String()

	runParams := []string{
		"run",
		"--detach",
		"--name",
		containerName,
		"--hostname",
		addContainer.Hostname.String(),
		"--env",
		"VIRTUAL_HOST=" + addContainer.Hostname.String(),
	}

	if len(addContainer.Envs) > 0 {
		envFlags := []string{}
		for _, env := range addContainer.Envs {
			envFlags = append(envFlags, "--env")
			envFlags = append(envFlags, env.String())
		}

		runParams = append(runParams, envFlags...)
	}

	baseSpecs, err := repo.getBaseSpecs(*addContainer.ResourceProfileId)
	if err != nil {
		return err
	}

	baseSpecsParams := []string{
		"--cpus",
		baseSpecs.CpuCores.String(),
		"--memory",
		baseSpecs.MemoryBytes.String(),
	}
	runParams = append(runParams, baseSpecsParams...)

	if addContainer.RestartPolicy != nil {
		runParams = append(runParams, "--restart", addContainer.RestartPolicy.String())
	}

	if addContainer.Entrypoint != nil {
		runParams = append(runParams, "--entrypoint", addContainer.Entrypoint.String())
	}

	if len(addContainer.PortBindings) > 0 {
		portBindingsParams := []string{}
		for _, portBinding := range addContainer.PortBindings {
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

	err = infraHelper.EnableLingering(addContainer.AccountId)
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
	currentContainer entity.Container,
	updateContainer dto.UpdateContainer,
) error {
	if updateContainer.Status != nil {
		shouldUpdateStatus := updateContainer.Status != &currentContainer.Status
		if !shouldUpdateStatus {
			return nil
		}

		actionToPerform := "start"
		if !*updateContainer.Status {
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

	if updateContainer.ResourceProfileId != nil {
		newSpecs, err := repo.getBaseSpecs(*updateContainer.ResourceProfileId)
		if err != nil {
			return err
		}

		_, err = infraHelper.RunCmdAsUser(
			updateContainer.AccountId,
			"podman",
			"update",
			"--cpus",
			newSpecs.CpuCores.String(),
			"--memory",
			newSpecs.MemoryBytes.String(),
			updateContainer.ContainerId.String(),
		)
		if err != nil {
			return err
		}
	}

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
