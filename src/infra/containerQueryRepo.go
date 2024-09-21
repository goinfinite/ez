package infra

import (
	"encoding/json"
	"errors"
	"log/slog"
	"strconv"
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

func (repo *ContainerQueryRepo) Read() ([]entity.Container, error) {
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
			slog.Debug(
				"ContainerModelToEntityError",
				slog.String("containerId", containerModel.ID),
				slog.Any("error", err),
			)
			continue
		}
		containers = append(containers, containerEntity)
	}

	return containers, nil
}

func (repo *ContainerQueryRepo) ReadById(
	containerId valueObject.ContainerId,
) (containerEntity entity.Container, err error) {
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

	containerEntity, err = containerModel.ToEntity()
	if err != nil {
		return containerEntity, err
	}

	return containerEntity, nil
}

func (repo *ContainerQueryRepo) ReadByHostname(
	hostname valueObject.Fqdn,
) (containerEntity entity.Container, err error) {
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

	containerEntity, err = containerModel.ToEntity()
	if err != nil {
		return containerEntity, err
	}

	return containerEntity, nil
}

func (repo *ContainerQueryRepo) ReadByAccountId(
	accountId valueObject.AccountId,
) ([]entity.Container, error) {
	containers := []entity.Container{}

	containerModels := []dbModel.Container{}
	err := repo.persistentDbSvc.Handler.
		Preload("PortBindings").
		Where("account_id = ?", accountId.Uint64()).
		Find(&containerModels).Error
	if err != nil {
		return containers, err
	}

	for _, containerModel := range containerModels {
		containerEntity, err := containerModel.ToEntity()
		if err != nil {
			slog.Debug(
				"ContainerModelToEntityError",
				slog.String("containerId", containerModel.ID),
				slog.Any("error", err),
			)
			continue
		}
		containers = append(containers, containerEntity)
	}

	return containers, nil
}

func (repo *ContainerQueryRepo) ReadByImageId(
	accountId valueObject.AccountId,
	imageId valueObject.ContainerImageId,
) ([]entity.Container, error) {
	containers := []entity.Container{}

	containerModels := []dbModel.Container{}
	err := repo.persistentDbSvc.Handler.
		Preload("PortBindings").
		Where("account_id = ? AND image_id = ?", accountId.Uint64(), imageId.String()).
		Find(&containerModels).Error
	if err != nil {
		return containers, err
	}

	for _, containerModel := range containerModels {
		containerEntity, err := containerModel.ToEntity()
		if err != nil {
			slog.Debug(
				"ContainerModelToEntityError",
				slog.String("containerId", containerModel.ID),
				slog.Any("error", err),
			)
			continue
		}
		containers = append(containers, containerEntity)
	}

	return containers, nil
}

func (repo *ContainerQueryRepo) containerMetricFactory(
	accountId valueObject.AccountId,
	containerMetricsJson string,
) (containerMetrics valueObject.ContainerMetrics, err error) {
	metricsInfo := map[string]interface{}{}
	err = json.Unmarshal([]byte(containerMetricsJson), &metricsInfo)
	if err != nil {
		return containerMetrics, err
	}

	rawContainerId, assertOk := metricsInfo["ContainerID"].(string)
	if !assertOk {
		return containerMetrics, errors.New("ContainerIdNotFound")
	}
	if len(rawContainerId) > 12 {
		rawContainerId = rawContainerId[:12]
	}
	containerId, err := valueObject.NewContainerId(rawContainerId)
	if err != nil {
		return containerMetrics, err
	}

	cpuPercent, assertOk := metricsInfo["CPU"].(float64)
	if !assertOk {
		return containerMetrics, errors.New("CpuPercentNotFound")
	}

	avgCpu, assertOk := metricsInfo["AvgCPU"].(float64)
	if !assertOk {
		return containerMetrics, errors.New("AvgCpuNotFound")
	}

	memPercent, assertOk := metricsInfo["MemPerc"].(float64)
	if !assertOk {
		return containerMetrics, errors.New("MemPercentNotFound")
	}

	rawMemBytes, assertOk := metricsInfo["MemUsage"].(float64)
	if !assertOk {
		return containerMetrics, errors.New("MemBytesNotFound")
	}
	memBytes, err := valueObject.NewByte(rawMemBytes)
	if err != nil {
		return containerMetrics, errors.New("MemBytesParseError")
	}

	rawBlockInput, assertOk := metricsInfo["BlockInput"].(float64)
	if !assertOk {
		return containerMetrics, errors.New("BlockInputNotFound")
	}
	blockInput, err := valueObject.NewByte(rawBlockInput)
	if err != nil {
		return containerMetrics, errors.New("BlockInputParseError")
	}

	rawBlockOutput, assertOk := metricsInfo["BlockOutput"].(float64)
	if !assertOk {
		return containerMetrics, errors.New("BlockOutputNotFound")
	}
	blockOutput, err := valueObject.NewByte(rawBlockOutput)
	if err != nil {
		return containerMetrics, errors.New("BlockOutputParseError")
	}

	rawNetInputTotal := 0.0
	rawNetOutputTotal := 0.0

	rawNetworks, assertOk := metricsInfo["Network"].(map[string]interface{})
	if !assertOk {
		return containerMetrics, errors.New("NetworkNotFound")
	}
	for _, rawNetwork := range rawNetworks {
		rawNetworkMap, assertOk := rawNetwork.(map[string]interface{})
		if !assertOk {
			continue
		}

		rawNetInput, assertOk := rawNetworkMap["RxBytes"].(float64)
		if !assertOk {
			continue
		}

		rawNetInputTotal += rawNetInput

		rawNetOutput, assertOk := rawNetworkMap["TxBytes"].(float64)
		if !assertOk {
			continue
		}

		rawNetOutputTotal += rawNetOutput
	}

	netInput, err := valueObject.NewByte(rawNetInputTotal)
	if err != nil {
		return containerMetrics, errors.New("NetInputParseError")
	}

	netOutput, err := valueObject.NewByte(rawNetOutputTotal)
	if err != nil {
		return containerMetrics, errors.New("NetOutputParseError")
	}

	// TODO: Implement cache so that we don't have to run this command every time.
	blockUsageStr, err := infraHelper.RunCmdAsUserWithSubShell(
		accountId,
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

	inodesCount, err := strconv.ParseUint(blockUsageParts[1], 10, 64)
	if err != nil {
		inodesCount = 0
	}

	return valueObject.NewContainerMetrics(
		containerId, infraHelper.RoundFloat(cpuPercent), avgCpu, memBytes,
		infraHelper.RoundFloat(memPercent), blockInput, blockOutput,
		blockBytes, inodesCount, netInput, netOutput,
	), nil
}

func (repo *ContainerQueryRepo) containerMetricsFactory(
	accountId valueObject.AccountId,
	containersMetricsStr string,
) (map[valueObject.ContainerId]valueObject.ContainerMetrics, error) {
	containersMetrics := map[valueObject.ContainerId]valueObject.ContainerMetrics{}
	if len(containersMetricsStr) == 0 {
		return containersMetrics, nil
	}

	containersMetricsList := strings.Split(containersMetricsStr, "\n")
	if len(containersMetricsList) == 0 {
		return containersMetrics, errors.New("ContainersMetricsParseError")
	}

	for metricsIndex, metricsJson := range containersMetricsList {
		containerMetrics, err := repo.containerMetricFactory(accountId, metricsJson)
		if err != nil {
			slog.Debug(
				"ContainerMetricsParseError",
				slog.Int("index", metricsIndex),
				slog.Any("error", err),
			)
			continue
		}

		containersMetrics[containerMetrics.ContainerId] = containerMetrics
	}

	return containersMetrics, nil
}

func (repo *ContainerQueryRepo) getWithMetricsByAccId(
	accountId valueObject.AccountId,
) ([]dto.ContainerWithMetrics, error) {
	containersWithMetrics := []dto.ContainerWithMetrics{}

	containersMetricsStr, err := infraHelper.RunCmdAsUser(
		accountId,
		"podman", "stats", "--no-stream", "--no-reset", "--format", "{{json .ContainerStats}}",
	)
	if err != nil {
		return containersWithMetrics, errors.New("AccPodmanStatsError" + err.Error())
	}

	runningContainersMetrics, err := repo.containerMetricsFactory(
		accountId, containersMetricsStr,
	)
	if err != nil {
		return containersWithMetrics, err
	}

	containerEntities, err := repo.ReadByAccountId(accountId)
	if err != nil {
		return containersWithMetrics, err
	}

	for _, containerEntity := range containerEntities {
		if _, exists := runningContainersMetrics[containerEntity.Id]; !exists {
			slog.Debug(
				"ContainerMetricsNotFound",
				slog.String("containerId", containerEntity.Id.String()),
			)
			continue
		}

		containerWithMetrics := dto.NewContainerWithMetrics(
			containerEntity, runningContainersMetrics[containerEntity.Id],
		)
		containersWithMetrics = append(containersWithMetrics, containerWithMetrics)
	}

	return containersWithMetrics, nil
}

func (repo *ContainerQueryRepo) ReadWithMetrics() ([]dto.ContainerWithMetrics, error) {
	containersWithMetrics := []dto.ContainerWithMetrics{}

	accountsList, err := NewAccountQueryRepo(repo.persistentDbSvc).Read()
	if err != nil {
		return containersWithMetrics, err
	}

	for _, account := range accountsList {
		accContainersWithMetrics, err := repo.getWithMetricsByAccId(account.Id)
		if err != nil {
			slog.Debug(
				"ReadAccountContainersWithMetricsError",
				slog.String("accountId", account.Id.String()),
				slog.Any("error", err),
			)
			continue
		}

		containersWithMetrics = append(containersWithMetrics, accContainersWithMetrics...)
	}

	return containersWithMetrics, nil
}

func (repo *ContainerQueryRepo) ReadWithMetricsById(
	containerId valueObject.ContainerId,
) (containerWithMetrics dto.ContainerWithMetrics, err error) {
	containerEntity, err := repo.ReadById(containerId)
	if err != nil {
		return containerWithMetrics, err
	}

	containersMetricStr, err := infraHelper.RunCmdAsUser(
		containerEntity.AccountId,
		"podman", "stats",
		"--no-stream", "--no-reset", "--format", "{{json .ContainerStats}}",
		containerId.String(),
	)
	if err != nil {
		return containerWithMetrics, errors.New("AccPodmanStatsError" + err.Error())
	}

	runningContainerMetrics, err := repo.containerMetricsFactory(
		containerEntity.AccountId, containersMetricStr,
	)
	if err != nil {
		return containerWithMetrics, err
	}
	if len(runningContainerMetrics) == 0 {
		return containerWithMetrics, errors.New("ContainerMetricsNotFound")
	}

	if _, exists := runningContainerMetrics[containerId]; !exists {
		return containerWithMetrics, errors.New("ContainerMetricsNotFound")
	}

	return dto.NewContainerWithMetrics(
		containerEntity, runningContainerMetrics[containerId],
	), nil
}
