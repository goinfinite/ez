package infra

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
	dbModel "github.com/speedianet/control/src/infra/db/model"
	infraHelper "github.com/speedianet/control/src/infra/helper"
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
		usedPrivatePorts := []valueObject.NetworkPort{}

		for pbIndex, portBindingVo := range addContainer.PortBindings {
			portBindingModel := dbModel.ContainerPortBinding{
				ContainerPort: uint(portBindingVo.ContainerPort),
				PublicPort:    uint(portBindingVo.PublicPort),
			}

			nextPrivatePort, err := portBindingModel.GetNextAvailablePrivatePort(
				repo.dbSvc.Orm,
				usedPrivatePorts,
			)
			if err != nil {
				return errors.New("FailedToGetNextAvailablePrivatePort")
			}

			usedPrivatePorts = append(usedPrivatePorts, nextPrivatePort)

			portBindingModel.PrivatePort = uint(nextPrivatePort.Get())
			addContainer.PortBindings[pbIndex].PrivatePort = &nextPrivatePort

			portBindingsParams = append(portBindingsParams, "--publish")
			portBindingsString := nextPrivatePort.String() +
				":" + portBindingVo.ContainerPort.String()

			protocolStr := portBindingVo.Protocol.String()
			if protocolStr != "" {
				portBindingModel.Protocol = protocolStr
				portBindingsString += "/" + protocolStr
			}

			portBindingsParams = append(portBindingsParams, portBindingsString)
		}

		runParams = append(runParams, portBindingsParams...)
	}

	runParams = append(runParams, addContainer.ImageAddr.String())

	err = infraHelper.EnableLingering(addContainer.AccountId)
	if err != nil {
		return errors.New("FailedToEnableLingering: " + err.Error())
	}
	time.Sleep(1 * time.Second)

	_, err = infraHelper.RunCmdAsUser(
		addContainer.AccountId,
		"podman",
		runParams...,
	)
	if err != nil {
		return errors.New("ContainerRunError: " + err.Error())
	}

	containerInfoJson, err := infraHelper.RunCmdAsUser(
		addContainer.AccountId,
		"podman",
		"container",
		"inspect",
		containerName,
		"--format",
		"{{json .}}",
	)
	if err != nil {
		return errors.New("GetContainerInfoError")
	}

	containerInfo := map[string]interface{}{}
	err = json.Unmarshal([]byte(containerInfoJson), &containerInfo)
	if err != nil {
		return errors.New("ContainerInfoParseError")
	}

	rawContainerId, assertOk := containerInfo["Id"].(string)
	if !assertOk {
		return errors.New("ContainerIdParseError")
	}
	containerId, err := valueObject.NewContainerId(rawContainerId)
	if err != nil {
		return err
	}

	rawImageHash, assertOk := containerInfo["ImageDigest"].(string)
	if !assertOk {
		return errors.New("ImageHashParseError")
	}
	rawImageHash = strings.TrimPrefix(rawImageHash, "sha256:")

	imageHash, err := valueObject.NewHash(rawImageHash)
	if err != nil {
		return err
	}

	nowUnixTime := valueObject.UnixTime(time.Now().Unix())

	containerEntity := entity.NewContainer(
		containerId,
		addContainer.AccountId,
		addContainer.Hostname,
		true,
		addContainer.ImageAddr,
		imageHash,
		addContainer.PortBindings,
		*addContainer.RestartPolicy,
		0,
		addContainer.Entrypoint,
		nowUnixTime,
		&nowUnixTime,
		*addContainer.ProfileId,
		addContainer.Envs,
	)

	containerModel := dbModel.Container{}.ToModel(containerEntity)

	createResult := repo.dbSvc.Orm.Create(&containerModel)
	if createResult.Error != nil {
		return createResult.Error
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

	if updateContainer.ProfileId == nil {
		return nil
	}

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
	return err
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
