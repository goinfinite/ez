package infra

import (
	"errors"
	"log"

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

	if addContainer.BaseSpecs != nil {
		specsParams := []string{
			"--cpus",
			addContainer.BaseSpecs.GetCpuCoresAsString(),
			"--memory",
			addContainer.BaseSpecs.GetMemoryAsString(),
			"--storage-opt",
			"size=" + addContainer.BaseSpecs.GetStorageAsString(),
		}

		runParams = append(runParams, specsParams...)
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
			portBindingsString := portBinding.GetHostPortAsString() + ":" + portBinding.GetContainerPortAsString()
			if portBinding.GetProtocol().String() != "" {
				portBindingsString += "/" + portBinding.GetProtocol().String()
			}
			portBindingsParams = append(portBindingsParams, portBindingsString)
		}

		runParams = append(runParams, portBindingsParams...)
	}

	runParams = append(runParams, addContainer.Image.String())

	_, err := infraHelper.RunCmd("podman", runParams...)
	if err != nil {
		log.Printf("AddContainerFailed: %v", err)
		return errors.New("AddContainerFailed")
	}

	return nil
}

func (repo ContainerCmdRepo) Update(dto.UpdateContainer) error {
	return nil
}

func (repo ContainerCmdRepo) Delete(valueObject.ContainerId) error {
	return nil
}
