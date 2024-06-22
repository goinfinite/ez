package infra

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
	dbModel "github.com/speedianet/control/src/infra/db/model"
	infraHelper "github.com/speedianet/control/src/infra/helper"
)

type ContainerCmdRepo struct {
	persistentDbSvc           *db.PersistentDatabaseService
	containerQueryRepo        *ContainerQueryRepo
	containerProfileQueryRepo *ContainerProfileQueryRepo
}

func NewContainerCmdRepo(persistentDbSvc *db.PersistentDatabaseService) *ContainerCmdRepo {
	containerQueryRepo := NewContainerQueryRepo(persistentDbSvc)
	ContainerProfileQueryRepo := NewContainerProfileQueryRepo(persistentDbSvc)
	return &ContainerCmdRepo{
		persistentDbSvc:           persistentDbSvc,
		containerQueryRepo:        containerQueryRepo,
		containerProfileQueryRepo: ContainerProfileQueryRepo,
	}
}

func (repo *ContainerCmdRepo) calibratePortBindings(
	originalPortBindings []valueObject.PortBinding,
) (calibratedPortBindings []valueObject.PortBinding, err error) {
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

		if portBindingVo.Protocol.String() == "udp" {
			portBindingsString += "/udp"
		}

		portBindingsParams = append(portBindingsParams, portBindingsString)
	}

	return portBindingsParams
}

func (repo *ContainerCmdRepo) containerEntityFactory(
	createDto dto.CreateContainer,
	containerName string,
) (containerEntity entity.Container, err error) {
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

func (repo *ContainerCmdRepo) getAccountHomeDir(
	accountId valueObject.AccountId,
) (string, error) {
	return infraHelper.RunCmdWithSubShell(
		"awk -F: '$3 == " + accountId.String() + " {print $6}' /etc/passwd",
	)
}

func (repo *ContainerCmdRepo) runLaunchScript(
	accountId valueObject.AccountId,
	containerId valueObject.ContainerId,
	launchScript *valueObject.LaunchScript,
) error {
	accountIdStr := accountId.String()
	accountHomeDir, err := repo.getAccountHomeDir(accountId)
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
) (containerId valueObject.ContainerId, err error) {
	hostnameStr := createDto.Hostname.String()
	containerName := hostnameStr

	createParams := []string{
		"create",
		"--cgroups=split",
		"--sdnotify=conmon",
		"--name", containerName,
		"--hostname", hostnameStr,
		"--env", "PRIMARY_VHOST=" + hostnameStr,
	}

	if len(createDto.Envs) > 0 {
		envFlags := []string{}
		for _, env := range createDto.Envs {
			envFlags = append(envFlags, "--env")
			envFlags = append(envFlags, env.String())
		}

		createParams = append(createParams, envFlags...)
	}

	if createDto.Entrypoint != nil {
		createParams = append(createParams, "--entrypoint", createDto.Entrypoint.String())
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

		createParams = append(createParams, portBindingsParams...)
	}

	createParams = append(createParams, createDto.ImageAddress.String())

	_, err = infraHelper.RunCmdAsUser(
		createDto.AccountId,
		"podman", createParams...,
	)
	if err != nil {
		return containerId, errors.New("CreateContainerError: " + err.Error())
	}

	containerProfile, err := repo.containerProfileQueryRepo.GetById(*createDto.ProfileId)
	if err != nil {
		return containerId, err
	}

	cpuQuotaPercentile := containerProfile.BaseSpecs.CpuCores.Get() * 100
	cpuQuotaPercentileStr := strconv.FormatFloat(cpuQuotaPercentile, 'f', -1, 64) + "%"

	containerEntity, err := repo.containerEntityFactory(createDto, containerName)
	if err != nil {
		return containerId, err
	}

	containerId = containerEntity.Id

	systemdUnitTemplate := `[Unit]
Description=` + containerName + ` Container
Wants=network-online.target
After=network-online.target
RequiresMountsFor=%t/containers

[Service]
Type=forking
Delegate=true
{{if .RestartPolicy}}Restart={{ .RestartPolicy.String }}{{end}}
Environment=PODMAN_SYSTEMD_UNIT=%n
CPUQuota=` + cpuQuotaPercentileStr + `
MemoryMax=` + containerProfile.BaseSpecs.MemoryBytes.String() + `
MemorySwapMax=0
ExecStart=/usr/bin/podman start ` + containerId.String() + `
ExecStop=/usr/bin/podman stop -t 10 ` + containerId.String() + `
TimeoutStartSec=900
TimeoutStopSec=60
OOMScoreAdjust=500
KillMode=mixed

[Install]
WantedBy=default.target
`

	systemdUnitTemplatePtr, err := template.New("systemdUnitFile").Parse(systemdUnitTemplate)
	if err != nil {
		return containerId, errors.New("SystemdUnitTemplateParsingError: " + err.Error())
	}

	var systemdUnitFileContent strings.Builder
	err = systemdUnitTemplatePtr.Execute(&systemdUnitFileContent, createDto)
	if err != nil {
		return containerId, errors.New("SystemdUnitTemplateExecutionError: " + err.Error())
	}

	accountHomeDir, err := repo.getAccountHomeDir(createDto.AccountId)
	if err != nil {
		return containerId, errors.New("GetAccountHomeDirError: " + err.Error())
	}

	systemdUnitDir := accountHomeDir + "/.config/systemd/user/"
	_, err = infraHelper.RunCmdAsUser(createDto.AccountId, "mkdir", "-p", systemdUnitDir)
	if err != nil {
		return containerId, errors.New("MakeSystemdUnitDirError: " + err.Error())
	}

	unitName := containerName + ".service"
	systemdUnitFilePath := systemdUnitDir + unitName
	err = infraHelper.UpdateFile(systemdUnitFilePath, systemdUnitFileContent.String(), true)
	if err != nil {
		return containerId, errors.New("WriteSystemdUnitFileError: " + err.Error())
	}

	accountIdStr := createDto.AccountId.String()
	_, err = infraHelper.RunCmd("chown", accountIdStr+":"+accountIdStr, systemdUnitFilePath)
	if err != nil {
		return containerId, errors.New("ChownSystemdUnitFileError: " + err.Error())
	}

	_, err = infraHelper.RunCmdAsUser(
		createDto.AccountId,
		"systemctl", "--user", "daemon-reload",
	)
	if err != nil {
		return containerId, errors.New("SystemdDaemonReloadError: " + err.Error())
	}

	_, err = infraHelper.RunCmdAsUser(
		createDto.AccountId,
		"systemctl", "--user", "enable", "--now", unitName,
	)
	if err != nil {
		return containerId, errors.New("SystemdEnableUnitError: " + err.Error())
	}

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

func (repo *ContainerCmdRepo) Update(updateDto dto.UpdateContainer) error {
	containerEntity, err := repo.containerQueryRepo.GetById(updateDto.ContainerId)
	if err != nil {
		return err
	}

	containerHostnameStr := containerEntity.Hostname.String()
	unitName := containerHostnameStr + ".service"
	containerIdStr := updateDto.ContainerId.String()
	containerModel := dbModel.Container{ID: containerIdStr}

	if updateDto.Status != nil && *updateDto.Status != containerEntity.Status {
		systemdCmd := "stop"
		if *updateDto.Status {
			systemdCmd = "start"
		}

		_, err = infraHelper.RunCmdAsUser(
			updateDto.AccountId,
			"systemctl", "--user", systemdCmd, unitName,
		)
		if err != nil {
			return errors.New("SystemdCmdError: " + err.Error())
		}

		err = repo.persistentDbSvc.Handler.
			Model(&containerModel).
			Update("status", *updateDto.Status).Error
		if err != nil {
			return err
		}
	}

	if updateDto.ProfileId == nil {
		return nil
	}

	containerProfile, err := repo.containerProfileQueryRepo.GetById(*updateDto.ProfileId)
	if err != nil {
		return err
	}

	cpuQuotaPercentile := containerProfile.BaseSpecs.CpuCores.Get() * 100
	cpuQuotaPercentileStr := strconv.FormatFloat(cpuQuotaPercentile, 'f', -1, 64) + "%"

	unitFilePath, err := infraHelper.RunCmdAsUser(
		updateDto.AccountId,
		"systemctl", "--user", "show", "-P", "FragmentPath", unitName,
	)
	if err != nil {
		return errors.New("GetSystemdUnitFilePathError: " + err.Error())
	}

	_, err = infraHelper.RunCmdAsUser(
		updateDto.AccountId,
		"sed", "-i", "s/CPUQuota=.*/CPUQuota="+cpuQuotaPercentileStr+"/",
		unitFilePath,
	)
	if err != nil {
		return errors.New("UpdateCpuQuotaError: " + err.Error())
	}

	_, err = infraHelper.RunCmdAsUser(
		updateDto.AccountId,
		"sed", "-i", "s/MemoryMax=.*/MemoryMax="+containerProfile.BaseSpecs.MemoryBytes.String()+"/",
		unitFilePath,
	)
	if err != nil {
		return errors.New("UpdateMemoryQuotaError: " + err.Error())
	}

	_, err = infraHelper.RunCmdAsUser(
		updateDto.AccountId,
		"systemctl", "--user", "daemon-reload",
	)
	if err != nil {
		return errors.New("SystemdDaemonReloadError: " + err.Error())
	}

	// Podman doesn't read the systemd unit file on reload, so it's necessary to
	// update the container specs directly as well.
	_, err = infraHelper.RunCmdAsUser(
		updateDto.AccountId,
		"podman", "update",
		"--cpus", containerProfile.BaseSpecs.CpuCores.String(),
		"--memory", containerProfile.BaseSpecs.MemoryBytes.String(),
		updateDto.ContainerId.String(),
	)
	if err != nil {
		ignorableError := "error opening file"
		if !strings.Contains(err.Error(), ignorableError) {
			return errors.New("UpdateContainerSpecsError: " + err.Error())
		}
	}

	return repo.persistentDbSvc.Handler.
		Model(&containerModel).
		Update("profile_id", updateDto.ProfileId.String()).Error
}

func (repo *ContainerCmdRepo) Delete(
	accountId valueObject.AccountId,
	containerId valueObject.ContainerId,
) error {
	containerEntity, err := repo.containerQueryRepo.GetById(containerId)
	if err != nil {
		return err
	}

	containerHostnameStr := containerEntity.Hostname.String()
	unitName := containerHostnameStr + ".service"

	_, err = infraHelper.RunCmdAsUser(
		accountId,
		"systemctl", "--user", "disable", "--now", unitName,
	)
	if err != nil {
		return errors.New("SystemdDisableUnitError: " + err.Error())
	}

	unitFilePath, err := infraHelper.RunCmdAsUser(
		accountId,
		"systemctl", "--user", "show", "-P", "FragmentPath", unitName,
	)
	if err != nil {
		return errors.New("GetSystemdUnitFilePathError: " + err.Error())
	}

	_, err = infraHelper.RunCmdAsUser(
		accountId,
		"rm", "-f", unitFilePath,
	)
	if err != nil {
		return errors.New("RemoveSystemdUnitFileError: " + err.Error())
	}

	_, err = infraHelper.RunCmdAsUser(
		accountId,
		"systemctl", "--user", "daemon-reload",
	)
	if err != nil {
		return errors.New("SystemdDaemonReloadError: " + err.Error())
	}

	containerIdStr := containerId.String()
	_, err = infraHelper.RunCmdAsUser(
		accountId,
		"podman", "rm", "--force", containerIdStr,
	)
	if err != nil {
		return errors.New("RemoveContainerError: " + err.Error())
	}

	portBindingModel := dbModel.ContainerPortBinding{}
	deleteResult := repo.persistentDbSvc.Handler.Delete(
		portBindingModel,
		"container_id = ?", containerIdStr,
	)
	if deleteResult.Error != nil {
		return err
	}

	containerModel := dbModel.Container{ID: containerIdStr}
	deleteResult = repo.persistentDbSvc.Handler.Delete(&containerModel)
	return deleteResult.Error
}

func (repo *ContainerCmdRepo) GenerateContainerSessionToken(
	autoLoginDto dto.ContainerAutoLogin,
) (tokenValue valueObject.AccessTokenValue, err error) {
	containerEntity, err := repo.containerQueryRepo.GetById(autoLoginDto.ContainerId)
	if err != nil {
		return tokenValue, errors.New("ContainerNotFound")
	}

	randomPassword := infraHelper.GenPass(16)
	_, _ = repo.runContainerCmd(
		containerEntity.AccountId, containerEntity.Id,
		"os account create -u control -p "+randomPassword,
	)

	_, err = repo.runContainerCmd(
		containerEntity.AccountId, containerEntity.Id,
		"os account update -u control -p "+randomPassword,
	)
	if err != nil {
		return tokenValue, errors.New("UpdateControlUserPasswordFailed: " + err.Error())
	}

	ipAddressStr := autoLoginDto.IpAddress.String()
	loginResponseJson, err := repo.runContainerCmd(
		containerEntity.AccountId, containerEntity.Id,
		"os auth login -u control -p "+randomPassword+" -i "+ipAddressStr,
	)
	if err != nil {
		return tokenValue, errors.New("LoginWithControlUserFailed: " + err.Error())
	}

	var loginResponseMap map[string]interface{}
	err = json.Unmarshal([]byte(loginResponseJson), &loginResponseMap)
	if err != nil {
		return tokenValue, errors.New("LoginResponseParseError: " + err.Error())
	}

	rawResponseBody, assertOk := loginResponseMap["body"].(map[string]interface{})
	if !assertOk || len(rawResponseBody) == 0 {
		return tokenValue, errors.New("LoginResponseBodyParseError")
	}

	rawTokenValue, assertOk := rawResponseBody["tokenStr"].(string)
	if !assertOk || len(rawTokenValue) == 0 {
		return tokenValue, errors.New("TokenValueParseError")
	}

	return valueObject.NewAccessTokenValue(rawTokenValue)
}
