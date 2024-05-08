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
	persistentDbSvc    *db.PersistentDatabaseService
	containerQueryRepo *ContainerQueryRepo
}

func NewContainerCmdRepo(persistentDbSvc *db.PersistentDatabaseService) *ContainerCmdRepo {
	containerQueryRepo := NewContainerQueryRepo(persistentDbSvc)
	return &ContainerCmdRepo{
		persistentDbSvc:    persistentDbSvc,
		containerQueryRepo: containerQueryRepo,
	}
}

func (repo *ContainerCmdRepo) getBaseSpecs(
	profileId valueObject.ContainerProfileId,
) (valueObject.ContainerSpecs, error) {
	profileQueryRepo := NewContainerProfileQueryRepo(repo.persistentDbSvc)
	containerProfile, err := profileQueryRepo.GetById(
		profileId,
	)
	if err != nil {
		return valueObject.ContainerSpecs{}, err
	}

	return containerProfile.BaseSpecs, nil
}

func (repo *ContainerCmdRepo) calibratePortBindings(
	originalPortBindings []valueObject.PortBinding,
) ([]valueObject.PortBinding, error) {
	calibratedPortBindings := []valueObject.PortBinding{}
	usedPrivatePorts := []valueObject.NetworkPort{}
	usedPublicPorts := []valueObject.NetworkPort{}
	portBindingModel := dbModel.ContainerPortBinding{}

	for _, originalPortBinding := range originalPortBindings {
		nextPrivatePort, err := portBindingModel.GetNextAvailablePrivatePort(
			repo.persistentDbSvc.Handler,
			usedPrivatePorts,
		)
		if err != nil {
			return calibratedPortBindings, errors.New(
				"GetNextPrivatePortError: + " + err.Error(),
			)
		}

		usedPrivatePorts = append(usedPrivatePorts, nextPrivatePort)

		calibratedPortBinding := valueObject.NewPortBinding(
			originalPortBinding.ServiceName,
			originalPortBinding.PublicPort,
			originalPortBinding.ContainerPort,
			originalPortBinding.Protocol,
			&nextPrivatePort,
		)

		if originalPortBinding.PublicPort.Get() == 0 {
			calibratedPortBindings = append(
				calibratedPortBindings,
				calibratedPortBinding,
			)
			continue
		}

		nextPublicPort, err := portBindingModel.GetNextAvailablePublicPort(
			repo.persistentDbSvc.Handler,
			calibratedPortBinding,
			usedPublicPorts,
		)
		if err != nil {
			return calibratedPortBindings, errors.New(
				"GetNextPublicPortError: " + err.Error(),
			)
		}

		usedPublicPorts = append(usedPublicPorts, nextPublicPort)

		calibratedPortBinding.PublicPort = nextPublicPort

		calibratedPortBindings = append(
			calibratedPortBindings,
			calibratedPortBinding,
		)
	}

	return calibratedPortBindings, nil
}

func (repo *ContainerCmdRepo) getPortBindingsParam(
	portBindings []valueObject.PortBinding,
) []string {
	portBindingsParams := []string{}
	for _, portBindingVo := range portBindings {
		portBindingsParams = append(portBindingsParams, "--publish")
		portBindingsString := portBindingVo.PrivatePort.String() +
			":" + portBindingVo.ContainerPort.String()

		protocolStr := portBindingVo.Protocol.String()
		if protocolStr != "" && protocolStr != "tcp" {
			portBindingsString += "/udp"
		}

		portBindingsParams = append(portBindingsParams, portBindingsString)
	}

	return portBindingsParams
}

func (repo *ContainerCmdRepo) containerEntityFactory(
	createDto dto.CreateContainer,
	containerName string,
) (entity.Container, error) {
	var containerEntity entity.Container

	containerInfoJson, err := infraHelper.RunCmdAsUser(
		createDto.AccountId,
		"podman", "container", "inspect", containerName, "--format", "{{json .}}",
	)
	if err != nil {
		return containerEntity, errors.New("GetContainerInfoError")
	}

	containerInfo := map[string]interface{}{}
	err = json.Unmarshal([]byte(containerInfoJson), &containerInfo)
	if err != nil {
		return containerEntity, errors.New("ContainerInfoParseError")
	}

	rawContainerId, assertOk := containerInfo["Id"].(string)
	if !assertOk || len(rawContainerId) < 12 {
		return containerEntity, errors.New("ContainerIdParseError")
	}

	rawContainerId = rawContainerId[:12]
	containerId, err := valueObject.NewContainerId(rawContainerId)
	if err != nil {
		return containerEntity, err
	}

	rawImageHash, assertOk := containerInfo["ImageDigest"].(string)
	if !assertOk {
		return containerEntity, errors.New("ImageHashParseError")
	}
	rawImageHash = strings.TrimPrefix(rawImageHash, "sha256:")

	imageHash, err := valueObject.NewHash(rawImageHash)
	if err != nil {
		return containerEntity, err
	}

	nowUnixTime := valueObject.UnixTime(time.Now().Unix())

	return entity.NewContainer(
		containerId,
		createDto.AccountId,
		createDto.Hostname,
		true,
		createDto.ImageAddress,
		imageHash,
		createDto.PortBindings,
		*createDto.RestartPolicy,
		0,
		createDto.Entrypoint,
		*createDto.ProfileId,
		createDto.Envs,
		nowUnixTime,
		nowUnixTime,
		&nowUnixTime,
		nil,
	), nil
}

func (repo *ContainerCmdRepo) runContainerCmd(
	accountId valueObject.AccountId,
	containerId valueObject.ContainerId,
	command string,
) (string, error) {
	return infraHelper.RunCmdAsUser(
		accountId,
		"podman", "exec", containerId.String(), "/bin/sh", "-c", command,
	)
}

func (repo *ContainerCmdRepo) runLaunchScript(
	accountId valueObject.AccountId,
	containerId valueObject.ContainerId,
	launchScript *valueObject.LaunchScript,
) error {
	accountIdStr := accountId.String()
	accountHomeDir, err := infraHelper.RunCmdWithSubShell(
		"awk -F: '$3 == " + accountIdStr + " {print $6}' /etc/passwd",
	)
	if err != nil {
		return errors.New("GetAccountHomeDirError: " + err.Error())
	}

	accountTmpDir := accountHomeDir + "/tmp"
	err = infraHelper.MakeDir(accountTmpDir)
	if err != nil {
		return errors.New("MakeTmpDirError: " + err.Error())
	}

	containerIdStr := containerId.String()
	launchScriptFilePath := accountTmpDir + "/launch-script-" + containerIdStr + ".sh"
	err = infraHelper.UpdateFile(launchScriptFilePath, launchScript.String(), true)
	if err != nil {
		return errors.New("WriteLaunchScriptError: " + err.Error())
	}

	_, err = infraHelper.RunCmd(
		"chown", "-R", accountIdStr+":"+accountIdStr, accountTmpDir,
	)
	if err != nil {
		return errors.New("ChownTmpDirError: " + err.Error())
	}

	_, err = infraHelper.RunCmdAsUser(
		accountId,
		"podman", "cp", launchScriptFilePath, containerIdStr+":/tmp/launch-script",
	)
	if err != nil {
		return errors.New("CopyLaunchScriptError: " + err.Error())
	}

	err = infraHelper.RemoveFile(launchScriptFilePath)
	if err != nil {
		return errors.New("RemoveLaunchScriptError: " + err.Error())
	}

	_, err = repo.runContainerCmd(
		accountId, containerId,
		"chmod +x /tmp/launch-script",
	)
	if err != nil {
		return errors.New("ChmodLaunchScriptError: " + err.Error())
	}

	_, err = repo.runContainerCmd(
		accountId, containerId,
		"/tmp/launch-script",
	)
	if err != nil {
		return errors.New("RunLaunchScriptError: " + err.Error())
	}

	_, err = repo.runContainerCmd(
		accountId, containerId,
		"rm -f /tmp/launch-script",
	)
	if err != nil {
		return errors.New("RemoveLaunchScriptError: " + err.Error())
	}

	return nil
}

func (repo *ContainerCmdRepo) Create(
	createDto dto.CreateContainer,
) (valueObject.ContainerId, error) {
	var containerId valueObject.ContainerId

	containerName := createDto.ProfileId.String() +
		"-" + createDto.Hostname.String()

	runParams := []string{
		"run", "--detach",
		"--name", containerName,
		"--hostname", createDto.Hostname.String(),
		"--env", "VIRTUAL_HOST=" + createDto.Hostname.String(),
	}

	if len(createDto.Envs) > 0 {
		envFlags := []string{}
		for _, env := range createDto.Envs {
			envFlags = append(envFlags, "--env")
			envFlags = append(envFlags, env.String())
		}

		runParams = append(runParams, envFlags...)
	}

	baseSpecs, err := repo.getBaseSpecs(*createDto.ProfileId)
	if err != nil {
		return containerId, err
	}

	baseSpecsParams := []string{
		"--cpus", baseSpecs.CpuCores.String(),
		"--memory", baseSpecs.MemoryBytes.String(),
	}
	runParams = append(runParams, baseSpecsParams...)

	if createDto.RestartPolicy != nil {
		runParams = append(runParams, "--restart", createDto.RestartPolicy.String())
	}

	if createDto.Entrypoint != nil {
		runParams = append(runParams, "--entrypoint", createDto.Entrypoint.String())
	}

	isSpeediaOs := createDto.ImageAddress.IsSpeediaOs()
	hasSpeediaOsPortBinding := false
	for _, portBinding := range createDto.PortBindings {
		if portBinding.ContainerPort.String() == "1618" {
			hasSpeediaOsPortBinding = true
			break
		}
	}

	if isSpeediaOs && !hasSpeediaOsPortBinding {
		portBinding, _ := valueObject.NewPortBindingFromString("sos")
		createDto.PortBindings = append(createDto.PortBindings, portBinding[0])
	}

	if len(createDto.PortBindings) > 0 {
		createDto.PortBindings, err = repo.calibratePortBindings(createDto.PortBindings)
		if err != nil {
			return containerId, err
		}

		portBindingsParams := repo.getPortBindingsParam(createDto.PortBindings)

		runParams = append(runParams, portBindingsParams...)
	}

	runParams = append(runParams, createDto.ImageAddress.String())

	_, err = infraHelper.RunCmdAsUser(
		createDto.AccountId,
		"podman", runParams...,
	)
	if err != nil {
		return containerId, errors.New("ContainerRunError: " + err.Error())
	}

	containerEntity, err := repo.containerEntityFactory(createDto, containerName)
	if err != nil {
		return containerId, err
	}

	containerId = containerEntity.Id

	containerModel := dbModel.Container{}.ToModel(containerEntity)
	createResult := repo.persistentDbSvc.Handler.Create(&containerModel)
	if createResult.Error != nil {
		return containerId, createResult.Error
	}

	if createDto.LaunchScript != nil {
		err = repo.runLaunchScript(
			createDto.AccountId, containerId, createDto.LaunchScript,
		)
		if err != nil {
			return containerId, errors.New("LaunchScriptError: " + err.Error())
		}
	}

	return containerId, err
}

func (repo *ContainerCmdRepo) updateContainerStatus(updateDto dto.UpdateContainer) error {
	actionToPerform := "start"
	if !*updateDto.Status {
		actionToPerform = "stop"
	}

	_, err := infraHelper.RunCmdAsUser(
		updateDto.AccountId,
		"podman", actionToPerform, updateDto.ContainerId.String(),
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

	updateResult := repo.persistentDbSvc.Handler.Model(&containerModel).Updates(updateMap)
	return updateResult.Error
}

func (repo *ContainerCmdRepo) Update(updateDto dto.UpdateContainer) error {
	currentContainer, err := repo.containerQueryRepo.GetById(updateDto.ContainerId)
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
		"podman", "update",
		"--cpus", newSpecs.CpuCores.String(),
		"--memory", newSpecs.MemoryBytes.String(),
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
		"podman", "rename", updateDto.ContainerId.String(), newContainerName,
	)
	if err != nil {
		return errors.New("FailedToRenameContainer: " + err.Error())
	}

	containerModel := dbModel.Container{ID: updateDto.ContainerId.String()}
	updateResult := repo.persistentDbSvc.Handler.Model(&containerModel).
		Update("profile_id", updateDto.ProfileId.String())
	return updateResult.Error
}

func (repo *ContainerCmdRepo) Delete(
	accId valueObject.AccountId,
	containerId valueObject.ContainerId,
) error {
	_, err := infraHelper.RunCmdAsUser(
		accId,
		"podman", "rm", "--force", containerId.String(),
	)
	if err != nil {
		return err
	}

	portBindingModel := dbModel.ContainerPortBinding{}
	deleteResult := repo.persistentDbSvc.Handler.Delete(
		portBindingModel,
		"container_id = ?",
		containerId.String(),
	)
	if deleteResult.Error != nil {
		return err
	}

	containerModel := dbModel.Container{ID: containerId.String()}
	deleteResult = repo.persistentDbSvc.Handler.Delete(&containerModel)
	return deleteResult.Error
}

func (repo *ContainerCmdRepo) GenerateContainerSessionToken(
	containerId valueObject.ContainerId,
) (valueObject.AccessTokenValue, error) {
	var tokenValue valueObject.AccessTokenValue

	containerEntity, err := repo.containerQueryRepo.GetById(containerId)
	if err != nil {
		return tokenValue, errors.New("ContainerNotFound")
	}

	randomPassword := infraHelper.GenPass(16)
	_, _ = repo.runContainerCmd(
		containerEntity.AccountId, containerId,
		"os account create -u control -p "+randomPassword,
	)

	_, err = repo.runContainerCmd(
		containerEntity.AccountId, containerId,
		"os account update -u control -p "+randomPassword,
	)
	if err != nil {
		return tokenValue, errors.New("UpdateControlUserPasswordFailed")
	}

	loginResponseJson, err := repo.runContainerCmd(
		containerEntity.AccountId, containerId,
		"os account login -u control -p "+randomPassword,
	)
	if err != nil {
		return tokenValue, errors.New("LoginWithControlUserFailed")
	}

	loginResponseMap := map[string]interface{}{}
	err = json.Unmarshal([]byte(loginResponseJson), &loginResponseMap)
	if err != nil {
		return tokenValue, errors.New("LoginResponseParseError")
	}

	rawResponseBody, assertOk := loginResponseMap["body"].(string)
	if !assertOk || len(rawResponseBody) == 0 {
		return tokenValue, errors.New("LoginResponseBodyParseError")
	}

	return valueObject.NewAccessTokenValue(rawResponseBody)
}
