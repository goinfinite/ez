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
	"gorm.io/gorm"
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

func (repo ContainerCmdRepo) Add(
	addDto dto.AddContainer,
) (valueObject.ContainerId, error) {
	var containerId valueObject.ContainerId

	containerName := addDto.ProfileId.String() +
		"-" + addDto.Hostname.String()

	runParams := []string{
		"run",
		"--detach",
		"--name",
		containerName,
		"--hostname",
		addDto.Hostname.String(),
		"--env",
		"VIRTUAL_HOST=" + addDto.Hostname.String(),
	}

	if len(addDto.Envs) > 0 {
		envFlags := []string{}
		for _, env := range addDto.Envs {
			envFlags = append(envFlags, "--env")
			envFlags = append(envFlags, env.String())
		}

		runParams = append(runParams, envFlags...)
	}

	baseSpecs, err := repo.getBaseSpecs(*addDto.ProfileId)
	if err != nil {
		return containerId, err
	}

	baseSpecsParams := []string{
		"--cpus",
		baseSpecs.CpuCores.String(),
		"--memory",
		baseSpecs.MemoryBytes.String(),
	}
	runParams = append(runParams, baseSpecsParams...)

	if addDto.RestartPolicy != nil {
		runParams = append(runParams, "--restart", addDto.RestartPolicy.String())
	}

	if addDto.Entrypoint != nil {
		runParams = append(runParams, "--entrypoint", addDto.Entrypoint.String())
	}

	if len(addDto.PortBindings) > 0 {
		portBindingsParams := []string{}
		usedPrivatePorts := []valueObject.NetworkPort{}

		for pbIndex, portBindingVo := range addDto.PortBindings {
			portBindingModel := dbModel.ContainerPortBinding{
				ContainerPort: uint(portBindingVo.ContainerPort),
				PublicPort:    uint(portBindingVo.PublicPort),
			}

			nextPrivatePort, err := portBindingModel.GetNextAvailablePrivatePort(
				repo.dbSvc.Orm,
				usedPrivatePorts,
			)
			if err != nil {
				return containerId, errors.New("FailedToGetNextAvailablePrivatePort")
			}

			usedPrivatePorts = append(usedPrivatePorts, nextPrivatePort)

			portBindingModel.PrivatePort = uint(nextPrivatePort.Get())
			addDto.PortBindings[pbIndex].PrivatePort = &nextPrivatePort

			portBindingsParams = append(portBindingsParams, "--publish")
			portBindingsString := nextPrivatePort.String() +
				":" + portBindingVo.ContainerPort.String()

			protocolStr := portBindingVo.Protocol.String()
			if protocolStr != "" && protocolStr != "tcp" {
				portBindingModel.Protocol = "udp"
				portBindingsString += "/udp"
			}

			portBindingsParams = append(portBindingsParams, portBindingsString)
		}

		runParams = append(runParams, portBindingsParams...)
	}

	runParams = append(runParams, addDto.ImageAddress.String())

	err = infraHelper.EnableLingering(addDto.AccountId)
	if err != nil {
		return containerId, errors.New("FailedToEnableLingering: " + err.Error())
	}
	time.Sleep(1 * time.Second)

	_, err = infraHelper.RunCmdAsUser(
		addDto.AccountId,
		"podman",
		runParams...,
	)
	if err != nil {
		return containerId, errors.New("ContainerRunError: " + err.Error())
	}

	containerInfoJson, err := infraHelper.RunCmdAsUser(
		addDto.AccountId,
		"podman",
		"container",
		"inspect",
		containerName,
		"--format",
		"{{json .}}",
	)
	if err != nil {
		return containerId, errors.New("GetContainerInfoError")
	}

	containerInfo := map[string]interface{}{}
	err = json.Unmarshal([]byte(containerInfoJson), &containerInfo)
	if err != nil {
		return containerId, errors.New("ContainerInfoParseError")
	}

	rawContainerId, assertOk := containerInfo["Id"].(string)
	if !assertOk || len(rawContainerId) < 12 {
		return containerId, errors.New("ContainerIdParseError")
	}

	rawContainerId = rawContainerId[:12]
	containerId, err = valueObject.NewContainerId(rawContainerId)
	if err != nil {
		return containerId, err
	}

	rawImageHash, assertOk := containerInfo["ImageDigest"].(string)
	if !assertOk {
		return containerId, errors.New("ImageHashParseError")
	}
	rawImageHash = strings.TrimPrefix(rawImageHash, "sha256:")

	imageHash, err := valueObject.NewHash(rawImageHash)
	if err != nil {
		return containerId, err
	}

	nowUnixTime := valueObject.UnixTime(time.Now().Unix())

	containerEntity := entity.NewContainer(
		containerId,
		addDto.AccountId,
		addDto.Hostname,
		true,
		addDto.ImageAddress,
		imageHash,
		addDto.PortBindings,
		*addDto.RestartPolicy,
		0,
		addDto.Entrypoint,
		*addDto.ProfileId,
		addDto.Envs,
		nowUnixTime,
		nowUnixTime,
		&nowUnixTime,
		nil,
	)

	containerModel := dbModel.Container{}.ToModel(containerEntity)

	createResult := repo.dbSvc.Orm.Create(&containerModel)
	if createResult.Error != nil {
		return containerId, createResult.Error
	}

	return containerId, nil
}

func (repo ContainerCmdRepo) updateContainerStatus(updateDto dto.UpdateContainer) error {
	actionToPerform := "start"
	if !*updateDto.Status {
		actionToPerform = "stop"
	}

	_, err := infraHelper.RunCmdAsUser(
		updateDto.AccountId,
		"podman",
		actionToPerform,
		updateDto.ContainerId.String(),
	)
	if err != nil {
		return err
	}

	containerModel := dbModel.Container{ID: updateDto.ContainerId.String()}
	updateMap := map[string]interface{}{
		"status":     *updateDto.Status,
		"started_at": time.Now(),
		"stopped_at": gorm.Expr("NULL"),
		"updated_at": time.Now(),
	}

	if !*updateDto.Status {
		updateMap["started_at"] = gorm.Expr("NULL")
		updateMap["stopped_at"] = time.Now()
	}

	updateResult := repo.dbSvc.Orm.Model(&containerModel).Updates(updateMap)
	return updateResult.Error
}

func (repo ContainerCmdRepo) Update(updateDto dto.UpdateContainer) error {
	containerQueryRepo := NewContainerQueryRepo(repo.dbSvc)
	currentContainer, err := containerQueryRepo.GetById(updateDto.ContainerId)
	if err != nil {
		return errors.New("ContainerNotFound")
	}

	if updateDto.Status != nil && *updateDto.Status != currentContainer.Status {
		err := repo.updateContainerStatus(updateDto)
		if err != nil {
			return errors.New("FailedToUpdateContainerStatus: " + err.Error())
		}

		// Current OCI implementations does not support permanent container resources
		// update. Therefore, when updating container status (on/off), it is also
		// necessary to update the container profile to guarantee that the container
		// resources are in accordance with the profile.
		updateDto.ProfileId = &currentContainer.ProfileId
	}

	if updateDto.ProfileId == nil {
		return nil
	}

	newSpecs, err := repo.getBaseSpecs(*updateDto.ProfileId)
	if err != nil {
		return err
	}

	_, err = infraHelper.RunCmdAsUser(
		updateDto.AccountId,
		"podman",
		"update",
		"--cpus",
		newSpecs.CpuCores.String(),
		"--memory",
		newSpecs.MemoryBytes.String(),
		updateDto.ContainerId.String(),
	)
	if err != nil {
		ignorableError := "error opening file"
		if !strings.Contains(err.Error(), ignorableError) {
			return errors.New("FailedToUpdateContainerSpecs: " + err.Error())
		}
	}

	newContainerName := updateDto.ProfileId.String() +
		"-" + currentContainer.Hostname.String()

	_, err = infraHelper.RunCmdAsUser(
		updateDto.AccountId,
		"podman",
		"rename",
		updateDto.ContainerId.String(),
		newContainerName,
	)
	if err != nil {
		return errors.New("FailedToRenameContainer: " + err.Error())
	}

	containerModel := dbModel.Container{ID: updateDto.ContainerId.String()}
	updateResult := repo.dbSvc.Orm.Model(&containerModel).
		Update("profile_id", updateDto.ProfileId.String())
	return updateResult.Error
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

	portBindingModel := dbModel.ContainerPortBinding{}
	deleteResult := repo.dbSvc.Orm.Delete(
		portBindingModel,
		"container_id = ?",
		containerId.String(),
	)
	if deleteResult.Error != nil {
		return err
	}

	containerModel := dbModel.Container{ID: containerId.String()}
	deleteResult = repo.dbSvc.Orm.Delete(&containerModel)
	return deleteResult.Error
}
