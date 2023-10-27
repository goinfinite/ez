package infra

import (
	"time"

	"github.com/goinfinite/fleet/src/domain/dto"
	"github.com/goinfinite/fleet/src/domain/entity"
	"github.com/goinfinite/fleet/src/domain/valueObject"
	"github.com/goinfinite/fleet/src/infra/db"
	infraHelper "github.com/goinfinite/fleet/src/infra/helper"
)

type ContainerCmdRepo struct {
	dbSvc *db.DatabaseService
}

func NewContainerCmdRepo(dbSvc *db.DatabaseService) *ContainerCmdRepo {
	return &ContainerCmdRepo{dbSvc: dbSvc}
}

func (repo ContainerCmdRepo) getBaseSpecs(
	profileId valueObject.ContainerProfileId,
) (valueObject.ContainerSpecs, error) {
	profileQueryRepo := NewContainerProfileQueryRepo(repo.dbSvc)
	containerProfile, err := profileQueryRepo.GetById(
		profileId,
	)
	if err != nil {
		return valueObject.ContainerSpecs{}, err
	}

	return containerProfile.BaseSpecs, nil
}

func (repo ContainerCmdRepo) Add(addContainer dto.AddContainer) error {
	containerName := addContainer.ProfileId.String() +
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

	baseSpecs, err := repo.getBaseSpecs(*addContainer.ProfileId)
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

func (repo ContainerCmdRepo) updateContainerStatus(
	currentContainer entity.Container,
	updateContainer dto.UpdateContainer,
) error {
	actionToPerform := "start"
	if !*updateContainer.Status {
		actionToPerform = "stop"
	}

	shouldUpdateStatus := updateContainer.Status != &currentContainer.Status
	if !shouldUpdateStatus {
		return nil
	}

	_, err := infraHelper.RunCmdAsUser(
		updateContainer.AccountId,
		"podman",
		actionToPerform,
		updateContainer.ContainerId.String(),
	)
	return err
}

func (repo ContainerCmdRepo) Update(
	currentContainer entity.Container,
	updateContainer dto.UpdateContainer,
) error {
	if updateContainer.Status != nil {
		err := repo.updateContainerStatus(currentContainer, updateContainer)
		if err != nil {
			return err
		}
	}

	if updateContainer.ProfileId != nil {
		newSpecs, err := repo.getBaseSpecs(*updateContainer.ProfileId)
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

		newContainerName := updateContainer.ProfileId.String() +
			"-" + currentContainer.Hostname.String()

		_, err = infraHelper.RunCmdAsUser(
			updateContainer.AccountId,
			"podman",
			"rename",
			updateContainer.ContainerId.String(),
			newContainerName,
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
