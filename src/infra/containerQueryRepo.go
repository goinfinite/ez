package infra

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
	dbModel "github.com/speedianet/control/src/infra/db/model"
	infraHelper "github.com/speedianet/control/src/infra/helper"
)

type ContainerQueryRepo struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewContainerQueryRepo(
	persistentDbSvc *db.PersistentDatabaseService,
) *ContainerQueryRepo {
	return &ContainerQueryRepo{persistentDbSvc: persistentDbSvc}
}

func (repo ContainerQueryRepo) Get() ([]entity.Container, error) {
	containers := []entity.Container{}

	containerModels := []dbModel.Container{}
	err := repo.persistentDbSvc.Handler.
		Preload("PortBindings").
		Find(&containerModels).Error
	if err != nil {
		return containers, err
	}

	for _, containerModel := range containerModels {
		containerEntity, err := containerModel.ToEntity()
		if err != nil {
			log.Printf("[%s] %s", containerModel.ID, err.Error())
			continue
		}
		containers = append(containers, containerEntity)
	}

	return containers, nil
}

func (repo ContainerQueryRepo) GetById(
	containerId valueObject.ContainerId,
) (entity.Container, error) {
	var containerEntity entity.Container

	var containerModel dbModel.Container
	queryResult := repo.persistentDbSvc.Handler.
		Preload("PortBindings").
		Where("id = ?", containerId.String()).
		Limit(1).
		Find(&containerModel)
	if queryResult.Error != nil {
		return containerEntity, queryResult.Error
	}

	if queryResult.RowsAffected == 0 {
		return containerEntity, errors.New("ContainerNotFound")
	}

	containerEntity, err := containerModel.ToEntity()
	if err != nil {
		return containerEntity, err
	}

	return containerEntity, nil
}

func (repo ContainerQueryRepo) GetByHostname(
	hostname valueObject.Fqdn,
) (entity.Container, error) {
	var containerEntity entity.Container

	var containerModel dbModel.Container
	queryResult := repo.persistentDbSvc.Handler.
		Preload("PortBindings").
		Where("hostname = ?", hostname.String()).
		Limit(1).
		Find(&containerModel)
	if queryResult.Error != nil {
		return containerEntity, queryResult.Error
	}

	if queryResult.RowsAffected == 0 {
		return containerEntity, errors.New("ContainerNotFound")
	}

	containerEntity, err := containerModel.ToEntity()
	if err != nil {
		return containerEntity, err
	}

	return containerEntity, nil
}

func (repo ContainerQueryRepo) GetByAccId(
	accId valueObject.AccountId,
) ([]entity.Container, error) {
	containers := []entity.Container{}

	containerModels := []dbModel.Container{}
	err := repo.persistentDbSvc.Handler.
		Preload("PortBindings").
		Where("account_id = ?", accId.Get()).
		Find(&containerModels).Error
	if err != nil {
		return containers, err
	}

	for _, containerModel := range containerModels {
		containerEntity, err := containerModel.ToEntity()
		if err != nil {
			log.Printf("[%s] %s", containerModel.ID, err.Error())
			continue
		}
		containers = append(containers, containerEntity)
	}

	return containers, nil
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

		cpuPercent, assertOk := containerMetricsInfo["CPU"].(float64)
		if !assertOk {
			continue
		}

		avgCpu, assertOk := containerMetricsInfo["AvgCPU"].(float64)
		if !assertOk {
			continue
		}

		memPercent, assertOk := containerMetricsInfo["MemPerc"].(float64)
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
			infraHelper.RoundFloat(cpuPercent),
			avgCpu,
			memBytes,
			infraHelper.RoundFloat(memPercent),
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

	accsList, err := NewAccQueryRepo(repo.persistentDbSvc).Get()
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
