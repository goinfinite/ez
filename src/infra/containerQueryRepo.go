package infra

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
	infraHelper "github.com/speedianet/control/src/infra/helper"
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

	rawRestartCount, assertOk := containerInfo["RestartCount"].(float64)
	if !assertOk {
		return entity.Container{}, errors.New("RestartCountParseError")
	}
	restartCount := uint64(rawRestartCount)

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
		restartCount,
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
		return []entity.Container{}, errors.New("AccPodmanListError" + err.Error())
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

func (repo ContainerQueryRepo) containerMetricsFactory(
	accountId valueObject.AccountId,
	containersMetricsStr string,
) (map[valueObject.ContainerId]valueObject.ContainerMetrics, error) {
	var containersMetrics = map[valueObject.ContainerId]valueObject.ContainerMetrics{}
	if len(containersMetricsStr) == 0 {
		return containersMetrics, nil
	}

	containersMetricsList := strings.Split(containersMetricsStr, "\n")
	if len(containersMetricsList) == 0 {
		return containersMetrics, errors.New("ContainersMetricsParseError")
	}

	for _, containerMetricsJson := range containersMetricsList {
		containerMetricsInfo := map[string]interface{}{}
		err := json.Unmarshal([]byte(containerMetricsJson), &containerMetricsInfo)
		if err != nil {
			continue
		}

		rawContainerId, assertOk := containerMetricsInfo["ContainerID"].(string)
		if !assertOk {
			continue
		}
		if len(rawContainerId) > 12 {
			rawContainerId = rawContainerId[:12]
		}
		containerId, err := valueObject.NewContainerId(rawContainerId)
		if err != nil {
			continue
		}

		cpuPerc, assertOk := containerMetricsInfo["CPU"].(float64)
		if !assertOk {
			continue
		}

		avgCpu, assertOk := containerMetricsInfo["AvgCPU"].(float64)
		if !assertOk {
			continue
		}

		memPerc, assertOk := containerMetricsInfo["MemPerc"].(float64)
		if !assertOk {
			continue
		}

		rawMemBytes, assertOk := containerMetricsInfo["MemUsage"].(float64)
		if !assertOk {
			continue
		}
		memBytes, err := valueObject.NewByte(rawMemBytes)
		if err != nil {
			continue
		}

		rawBlockInput, assertOk := containerMetricsInfo["BlockInput"].(float64)
		if !assertOk {
			continue
		}
		blockInput, err := valueObject.NewByte(rawBlockInput)
		if err != nil {
			continue
		}

		rawBlockOutput, assertOk := containerMetricsInfo["BlockOutput"].(float64)
		if !assertOk {
			continue
		}
		blockOutput, err := valueObject.NewByte(rawBlockOutput)
		if err != nil {
			continue
		}

		rawNetInput, assertOk := containerMetricsInfo["NetInput"].(float64)
		if !assertOk {
			continue
		}
		netInput, err := valueObject.NewByte(rawNetInput)
		if err != nil {
			continue
		}

		rawNetOutput, assertOk := containerMetricsInfo["NetOutput"].(float64)
		if !assertOk {
			continue
		}
		netOutput, err := valueObject.NewByte(rawNetOutput)
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

		containerMetrics := valueObject.NewContainerMetrics(
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

		containersMetrics[containerId] = containerMetrics
	}

	return containersMetrics, nil
}

func (repo ContainerQueryRepo) getWithMetricsByAccId(
	accId valueObject.AccountId,
) ([]dto.ContainerWithMetrics, error) {
	var containersWithMetrics []dto.ContainerWithMetrics

	containersMetricsStr, err := infraHelper.RunCmdAsUser(
		accId,
		"podman",
		"stats",
		"--no-stream",
		"--no-reset",
		"--format",
		"{{json .ContainerStats}}",
	)
	if err != nil {
		return containersWithMetrics, errors.New("AccPodmanStatsError" + err.Error())
	}

	runningContainersMetrics, err := repo.containerMetricsFactory(
		accId,
		containersMetricsStr,
	)
	if err != nil {
		return containersWithMetrics, err
	}

	containerEntities, err := repo.GetByAccId(accId)
	if err != nil {
		return containersWithMetrics, err
	}

	for _, container := range containerEntities {
		containerMetrics := valueObject.NewContainerMetricsWithBlankValues()

		for runningContainerId, runningContainerMetrics := range runningContainersMetrics {
			if runningContainerId != container.Id {
				continue
			}
			containerMetrics = runningContainerMetrics
		}

		containerWithMetrics := dto.NewContainerWithMetrics(
			container,
			containerMetrics,
		)
		containersWithMetrics = append(containersWithMetrics, containerWithMetrics)
	}

	return containersWithMetrics, nil
}

func (repo ContainerQueryRepo) GetWithMetrics() ([]dto.ContainerWithMetrics, error) {
	containersWithMetrics := []dto.ContainerWithMetrics{}

	accsList, err := NewAccQueryRepo(repo.dbSvc).Get()
	if err != nil {
		return containersWithMetrics, err
	}

	for _, acc := range accsList {
		accContainersWithMetrics, err := repo.getWithMetricsByAccId(acc.Id)
		if err != nil {
			log.Printf("AccId '%s' skipped: %s", acc.Id.String(), err.Error())
			continue
		}

		containersWithMetrics = append(containersWithMetrics, accContainersWithMetrics...)
	}

	return containersWithMetrics, nil
}
