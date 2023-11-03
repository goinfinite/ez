package infra

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/goinfinite/fleet/src/domain/dto"
	"github.com/goinfinite/fleet/src/domain/entity"
	"github.com/goinfinite/fleet/src/domain/valueObject"
	"github.com/goinfinite/fleet/src/infra/db"
	infraHelper "github.com/goinfinite/fleet/src/infra/helper"
)

type ContainerQueryRepo struct {
	dbSvc *db.DatabaseService
}

func NewContainerQueryRepo(dbSvc *db.DatabaseService) *ContainerQueryRepo {
	return &ContainerQueryRepo{dbSvc: dbSvc}
}

func (repo ContainerQueryRepo) parsePortBindings(
	rawPortBindings map[string]interface{},
) ([]valueObject.PortBinding, error) {
	portBindings := []valueObject.PortBinding{}

	for containerPortProtocol, hostIpHostPortsMap := range rawPortBindings {
		containerPortProtocol := strings.Split(containerPortProtocol, "/")
		if len(containerPortProtocol) != 2 {
			return portBindings, errors.New("PortBindingInfoParseError")
		}

		containerPort, err := valueObject.NewNetworkPort(containerPortProtocol[0])
		if err != nil {
			return portBindings, err
		}

		networkProtocol, err := valueObject.NewNetworkProtocol(containerPortProtocol[1])
		if err != nil {
			return portBindings, err
		}

		for _, hostIpHostPort := range hostIpHostPortsMap.([]interface{}) {
			hostIpHostPort, assertOk := hostIpHostPort.(map[string]interface{})
			if !assertOk {
				return portBindings, errors.New("PortBindingInfoParseError")
			}

			rawHostPort, assertOk := hostIpHostPort["HostPort"].(string)
			if !assertOk {
				return portBindings, errors.New("HostPortParseError")
			}

			networkPort, err := valueObject.NewNetworkPort(rawHostPort)
			if err != nil {
				return portBindings, err
			}

			portBinding := valueObject.NewPortBinding(
				networkProtocol,
				containerPort,
				networkPort,
			)
			portBindings = append(portBindings, portBinding)
		}
	}

	return portBindings, nil
}

func (repo ContainerQueryRepo) GetById(
	accId valueObject.AccountId,
	containerId valueObject.ContainerId,
) (entity.Container, error) {
	containerInfoJson, err := infraHelper.RunCmdAsUser(
		accId,
		"podman",
		"container",
		"inspect",
		containerId.String(),
		"--format",
		"{{json .}}",
	)
	if err != nil {
		return entity.Container{}, err
	}

	containerInfo := map[string]interface{}{}
	err = json.Unmarshal([]byte(containerInfoJson), &containerInfo)
	if err != nil {
		return entity.Container{}, err
	}

	rawConfig, assertOk := containerInfo["Config"].(map[string]interface{})
	if !assertOk {
		return entity.Container{}, errors.New("PodmanConfigParseError")
	}

	rawHostname, assertOk := rawConfig["Hostname"].(string)
	if !assertOk {
		return entity.Container{}, errors.New("HostnameParseError")
	}
	hostname, err := valueObject.NewFqdn(rawHostname)
	if err != nil {
		return entity.Container{}, err
	}

	rawState, assertOk := containerInfo["State"].(map[string]interface{})
	if !assertOk {
		return entity.Container{}, errors.New("StateParseError")
	}

	status, assertOk := rawState["Running"].(bool)
	if !assertOk {
		return entity.Container{}, errors.New("StatusParseError")
	}

	rawImage, assertOk := containerInfo["ImageName"].(string)
	if !assertOk {
		return entity.Container{}, errors.New("ImageParseError")
	}
	image, err := valueObject.NewContainerImgAddress(rawImage)
	if err != nil {
		return entity.Container{}, err
	}

	rawImageHash, assertOk := containerInfo["ImageDigest"].(string)
	if !assertOk {
		return entity.Container{}, errors.New("ImageHashParseError")
	}
	rawImageHash = strings.TrimLeft(rawImageHash, "sha256:")
	imageHash, err := valueObject.NewHash(rawImageHash)
	if err != nil {
		return entity.Container{}, err
	}

	rawNetworkSettings, assertOk := containerInfo["NetworkSettings"].(map[string]interface{})
	if !assertOk {
		return entity.Container{}, errors.New("NetworkSettingsParseError")
	}

	rawPrivateIpAddress, assertOk := rawNetworkSettings["IPAddress"].(string)
	if !assertOk {
		return entity.Container{}, errors.New("PrivateIpAddressParseError")
	}
	privateIpAddress, err := valueObject.NewIpAddress(rawPrivateIpAddress)
	if err != nil {
		privateIpAddress, _ = valueObject.NewIpAddress("0.0.0.0")
	}

	rawHostConfig, assertOk := containerInfo["HostConfig"].(map[string]interface{})
	if !assertOk {
		return entity.Container{}, errors.New("HostConfigParseError")
	}

	rawPortBindings, assertOk := rawHostConfig["PortBindings"].(map[string]interface{})
	if !assertOk {
		return entity.Container{}, errors.New("PortBindingsParseError")
	}
	portBindings, err := repo.parsePortBindings(rawPortBindings)
	if err != nil {
		return entity.Container{}, err
	}

	rawRestartPolicy, assertOk := rawHostConfig["RestartPolicy"].(map[string]interface{})
	if !assertOk {
		return entity.Container{}, errors.New("RestartPolicyParseError")
	}
	rawRestartPolicyName := rawRestartPolicy["Name"].(string)
	if rawRestartPolicyName == "" {
		rawRestartPolicyName = "no"
	}
	restartPolicyName, err := valueObject.NewContainerRestartPolicy(rawRestartPolicyName)
	if err != nil {
		return entity.Container{}, err
	}

	rawEntryPoint, assertOk := rawConfig["Entrypoint"].(string)
	if !assertOk {
		return entity.Container{}, errors.New("EntrypointParseError")
	}
	entrypoint, err := valueObject.NewContainerEntrypoint(rawEntryPoint)
	if err != nil {
		return entity.Container{}, err
	}

	rawCreatedAt, assertOk := containerInfo["Created"].(string)
	if !assertOk {
		return entity.Container{}, errors.New("CreatedAtParseError")
	}
	createdAtTime, err := time.Parse(time.RFC3339, rawCreatedAt)
	if err != nil {
		return entity.Container{}, err
	}
	createdAt := valueObject.UnixTime(createdAtTime.Unix())

	var startedAtPtr *valueObject.UnixTime
	rawStartedAt, assertOk := rawState["StartedAt"].(string)
	if !assertOk {
		return entity.Container{}, errors.New("StartedAtParseError")
	}
	startedAtTime, err := time.Parse(time.RFC3339, rawStartedAt)
	if err == nil {
		startedAt := valueObject.UnixTime(startedAtTime.Unix())
		startedAtPtr = &startedAt
	}

	containerName, assertOk := containerInfo["Name"].(string)
	if !assertOk {
		return entity.Container{}, errors.New("ContainerNameNotString")
	}

	containerNameParts := strings.Split(containerName, "-")
	if len(containerNameParts) < 2 {
		return entity.Container{}, errors.New("ContainerNameParseError")
	}

	profileId, err := valueObject.NewContainerProfileId(containerNameParts[0])
	if err != nil {
		return entity.Container{}, err
	}

	rawEnvs, assertOk := rawConfig["Env"].([]interface{})
	if !assertOk {
		return entity.Container{}, errors.New("EnvParseError")
	}
	envs := []valueObject.ContainerEnv{}
	for _, rawEnv := range rawEnvs {
		env, err := valueObject.NewContainerEnv(rawEnv.(string))
		if err != nil {
			continue
		}
		envs = append(envs, env)
	}

	return entity.NewContainer(
		containerId,
		accId,
		hostname,
		status,
		image,
		imageHash,
		privateIpAddress,
		portBindings,
		restartPolicyName,
		entrypoint,
		createdAt,
		startedAtPtr,
		profileId,
		envs,
	), nil
}

func (repo ContainerQueryRepo) GetByAccId(
	accId valueObject.AccountId,
) ([]entity.Container, error) {
	containersIds, err := infraHelper.RunCmdAsUser(
		accId,
		"podman",
		"container",
		"list",
		"--all",
		"--format",
		"{{.ID}}",
	)
	if err != nil {
		return []entity.Container{}, err
	}
	containersIdsList := strings.Split(containersIds, "\n")
	if len(containersIdsList) == 0 {
		return []entity.Container{}, nil
	}

	containers := []entity.Container{}
	for _, containerIdStr := range containersIdsList {
		containerIdStr = strings.TrimSpace(containerIdStr)
		containerId, err := valueObject.NewContainerId(containerIdStr)
		if err != nil {
			continue
		}

		container, err := repo.GetById(accId, containerId)
		if err != nil {
			log.Printf(
				"ContainerId '%s' skipped. ParseError: %s",
				containerId.String(),
				err.Error(),
			)
			continue
		}
		containers = append(containers, container)
	}

	return containers, nil
}

func (repo ContainerQueryRepo) Get() ([]entity.Container, error) {
	allContainers := []entity.Container{}

	accsList, err := NewAccQueryRepo(repo.dbSvc).Get()
	if err != nil {
		return allContainers, err
	}

	for _, acc := range accsList {
		accContainers, err := repo.GetByAccId(acc.Id)
		if err != nil {
			log.Printf(
				"AccId '%s' skipped. ParseError: %s",
				acc.Id.String(),
				err.Error(),
			)
			continue
		}
		allContainers = append(allContainers, accContainers...)
	}

	return allContainers, nil
}

func (repo ContainerQueryRepo) containerResourceUsageFactory(
	accountId valueObject.AccountId,
	containersUsageStr string,
) ([]valueObject.ContainerResourceUsage, error) {
	var containersUsage []valueObject.ContainerResourceUsage
	if len(containersUsageStr) == 0 {
		return containersUsage, nil
	}

	containersUsageList := strings.Split(containersUsageStr, "\n")
	if len(containersUsageList) == 0 {
		return containersUsage, errors.New("ContainersUsageParseError")
	}

	for _, containerUsageStr := range containersUsageList {
		containerUsageStr = strings.TrimSpace(containerUsageStr)
		containerUsageParts := strings.Split(containerUsageStr, " ")
		if len(containerUsageParts) != 9 {
			continue
		}

		containerId, err := valueObject.NewContainerId(containerUsageParts[0])
		if err != nil {
			continue
		}

		cpuPerc, err := strconv.ParseFloat(containerUsageParts[1], 64)
		if err != nil {
			continue
		}

		avgCpu, err := strconv.ParseFloat(containerUsageParts[2], 64)
		if err != nil {
			continue
		}

		memPerc, err := strconv.ParseFloat(containerUsageParts[3], 64)
		if err != nil {
			continue
		}

		memBytes, err := valueObject.NewByte(containerUsageParts[4])
		if err != nil {
			continue
		}

		blockInput, err := valueObject.NewByte(containerUsageParts[5])
		if err != nil {
			continue
		}

		blockOutput, err := valueObject.NewByte(containerUsageParts[6])
		if err != nil {
			continue
		}

		netInput, err := valueObject.NewByte(containerUsageParts[7])
		if err != nil {
			continue
		}

		netOutput, err := valueObject.NewByte(containerUsageParts[8])
		if err != nil {
			continue
		}

		blockUsageStr, err := infraHelper.RunCmdAsUser(
			accountId,
			"bash",
			"-c",
			"timeout 1 podman exec -it "+containerId.String()+
				" df --output=used,iused / | tail -n 1",
		)
		if err != nil {
			blockUsageStr = "0 0"
		}
		blockUsageStr = strings.TrimSpace(blockUsageStr)
		blockUsageParts := strings.Split(blockUsageStr, " ")
		if len(blockUsageParts) != 2 {
			blockUsageParts = []string{"0", "0"}
		}

		blockBytes, err := valueObject.NewByte(blockUsageParts[0])
		if err != nil {
			blockBytes, _ = valueObject.NewByte(0)
		}

		inodesCount, err := valueObject.NewInodesCount(blockUsageParts[1])
		if err != nil {
			inodesCount, _ = valueObject.NewInodesCount(0)
		}

		containerUsage := valueObject.NewContainerResourceUsage(
			accountId,
			containerId,
			cpuPerc,
			avgCpu,
			memBytes,
			memPerc,
			blockInput,
			blockOutput,
			blockBytes,
			inodesCount,
			netInput,
			netOutput,
		)

		containersUsage = append(containersUsage, containerUsage)
	}

	return containersUsage, nil
}

func (repo ContainerQueryRepo) getWithUsageByAccId(
	accId valueObject.AccountId,
) ([]dto.ContainerWithUsage, error) {
	var containersWithUsage []dto.ContainerWithUsage

	containersUsageStr, err := infraHelper.RunCmdAsUser(
		accId,
		"podman",
		"stats",
		"--no-stream",
		"--no-reset",
		"--format",
		"{{.ID}} {{.ContainerStats.CPU}} {{.ContainerStats.AvgCPU}} "+
			"{{.ContainerStats.MemPerc}} {{.ContainerStats.MemUsage}} "+
			"{{.ContainerStats.BlockInput}} {{.ContainerStats.BlockOutput}} "+
			"{{.ContainerStats.NetInput}} {{.ContainerStats.NetOutput}}",
	)
	if err != nil {
		return containersWithUsage, errors.New("AccPodmanStatsError" + err.Error())
	}

	containersUsage, err := repo.containerResourceUsageFactory(
		accId,
		containersUsageStr,
	)
	if err != nil {
		return containersWithUsage, err
	}

	containerEntities, err := repo.GetByAccId(accId)
	if err != nil {
		return containersWithUsage, err
	}

	for _, container := range containerEntities {
		for _, containerUsage := range containersUsage {
			if containerUsage.ContainerId != container.Id {
				continue
			}

			containerWithUsage := dto.NewContainerWithUsage(
				container,
				containerUsage,
			)
			containersWithUsage = append(containersWithUsage, containerWithUsage)
		}
	}

	return containersWithUsage, nil
}

func (repo ContainerQueryRepo) GetWithUsage() ([]dto.ContainerWithUsage, error) {
	var containersWithUsage []dto.ContainerWithUsage

	accsList, err := NewAccQueryRepo(repo.dbSvc).Get()
	if err != nil {
		return containersWithUsage, err
	}

	for _, acc := range accsList {
		accContainersWithUsage, err := repo.getWithUsageByAccId(acc.Id)
		if err != nil {
			log.Printf("AccId '%s' skipped: %s", acc.Id.String(), err.Error())
			continue
		}

		containersWithUsage = append(containersWithUsage, accContainersWithUsage...)
	}

	return containersWithUsage, nil
}
