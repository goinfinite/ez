package infra

import (
	"encoding/json"
	"errors"
	"log/slog"
	"math"
	"strconv"
	"strings"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
	"github.com/goinfinite/ez/src/infra/db"
	dbModel "github.com/goinfinite/ez/src/infra/db/model"
	infraHelper "github.com/goinfinite/ez/src/infra/helper"
	"github.com/iancoleman/strcase"
)

type ContainerQueryRepo struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewContainerQueryRepo(
	persistentDbSvc *db.PersistentDatabaseService,
) *ContainerQueryRepo {
	return &ContainerQueryRepo{persistentDbSvc: persistentDbSvc}
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
		containerId, math.Round(cpuPercent), avgCpu, memBytes,
		math.Round(memPercent), blockInput, blockOutput,
		blockBytes, inodesCount, netInput, netOutput,
	), nil
}

func (repo *ContainerQueryRepo) containerMetricsFactory(
	accountId valueObject.AccountId,
	rawContainersMetrics string,
) ([]valueObject.ContainerMetrics, error) {
	containersMetrics := []valueObject.ContainerMetrics{}
	if len(rawContainersMetrics) == 0 {
		return containersMetrics, nil
	}

	rawContainersMetricsEntries := strings.Split(rawContainersMetrics, "\n")
	if len(rawContainersMetricsEntries) == 0 {
		return containersMetrics, errors.New("ContainersMetricsEntriesNotFound")
	}

	for metricsIndex, metricsJson := range rawContainersMetricsEntries {
		containerMetrics, err := repo.containerMetricFactory(accountId, metricsJson)
		if err != nil {
			slog.Debug(
				"ContainerMetricsParseError",
				slog.Int("index", metricsIndex),
				slog.Any("error", err),
			)
			continue
		}

		containersMetrics = append(containersMetrics, containerMetrics)
	}

	return containersMetrics, nil
}

func (repo *ContainerQueryRepo) containersWithMetricsFactory(
	containerEntities []entity.Container,
) ([]dto.ContainerWithMetrics, error) {
	containersWithMetrics := []dto.ContainerWithMetrics{}

	uniqueAccountIdsMap := map[valueObject.AccountId]struct{}{}
	for _, containerEntity := range containerEntities {
		if _, exists := uniqueAccountIdsMap[containerEntity.AccountId]; exists {
			continue
		}

		uniqueAccountIdsMap[containerEntity.AccountId] = struct{}{}
	}

	// TODO: Add coroutine for parallel execution.
	runningContainersMetrics := []valueObject.ContainerMetrics{}
	for accountId := range uniqueAccountIdsMap {
		containersMetricsStr, err := infraHelper.RunCmdAsUser(
			accountId,
			"podman", "stats", "--no-stream", "--no-reset", "--format", "{{json .ContainerStats}}",
		)
		if err != nil {
			slog.Debug(
				"AccountPodmanStatsError",
				slog.Uint64("accountId", accountId.Uint64()),
				slog.Any("error", err),
			)
			continue
		}

		accountContainersMetrics, err := repo.containerMetricsFactory(
			accountId, containersMetricsStr,
		)
		if err != nil {
			slog.Debug(
				"AccountContainersMetricsParseError",
				slog.Uint64("accountId", accountId.Uint64()),
				slog.Any("error", err),
			)
			continue
		}

		runningContainersMetrics = append(runningContainersMetrics, accountContainersMetrics...)
	}

	containerIdMetricsMap := map[valueObject.ContainerId]valueObject.ContainerMetrics{}
	for _, containerMetrics := range runningContainersMetrics {
		containerIdMetricsMap[containerMetrics.ContainerId] = containerMetrics
	}

	for _, containerEntity := range containerEntities {
		containerMetrics := valueObject.NewBlankContainerMetrics(containerEntity.Id)
		if _, exists := containerIdMetricsMap[containerEntity.Id]; exists {
			containerMetrics = containerIdMetricsMap[containerEntity.Id]
		}

		containerWithMetrics := dto.NewContainerWithMetrics(
			containerEntity, containerMetrics,
		)
		containersWithMetrics = append(containersWithMetrics, containerWithMetrics)
	}

	return containersWithMetrics, nil
}

func (repo *ContainerQueryRepo) Read(
	requestDto dto.ReadContainersRequest,
) (responseDto dto.ReadContainersResponse, err error) {
	containerModel := dbModel.Container{}
	if requestDto.ContainerId != nil {
		containerModel.ID = requestDto.ContainerId.String()
	}
	if requestDto.ContainerAccountId != nil {
		containerModel.AccountID = requestDto.ContainerAccountId.Uint64()
	}
	if requestDto.ContainerHostname != nil {
		containerModel.Hostname = requestDto.ContainerHostname.String()
	}
	if requestDto.ContainerStatus != nil {
		containerModel.Status = *requestDto.ContainerStatus
	}
	if requestDto.ContainerImageId != nil {
		containerModel.ImageId = requestDto.ContainerImageId.String()
	}
	if requestDto.ContainerImageAddress != nil {
		containerModel.ImageAddress = requestDto.ContainerImageAddress.String()
	}
	if requestDto.ContainerImageHash != nil {
		containerModel.ImageHash = requestDto.ContainerImageHash.String()
	}
	for _, portBinding := range requestDto.ContainerPortBindings {
		portBindingModel := dbModel.ContainerPortBinding{
			ServiceName:   portBinding.ServiceName.String(),
			PublicPort:    portBinding.PublicPort.Uint16(),
			ContainerPort: portBinding.ContainerPort.Uint16(),
			Protocol:      portBinding.Protocol.String(),
		}
		if portBinding.PrivatePort != nil {
			portBindingModel.PrivatePort = portBinding.PrivatePort.Uint16()
		}

		containerModel.PortBindings = append(containerModel.PortBindings, portBindingModel)
	}
	if requestDto.ContainerRestartPolicy != nil {
		containerModel.RestartPolicy = requestDto.ContainerRestartPolicy.String()
	}
	if requestDto.ContainerProfileId != nil {
		containerModel.ProfileID = requestDto.ContainerProfileId.Uint64()
	}
	if len(requestDto.ContainerEnv) > 0 {
		envs := dbModel.Container{}.JoinEnvs(requestDto.ContainerEnv)
		containerModel.Envs = &envs
	}

	dbQuery := repo.persistentDbSvc.Handler.
		Model(&containerModel).
		Where(&containerModel).
		Preload("PortBindings")

	if requestDto.CreatedBeforeAt != nil {
		dbQuery = dbQuery.Where("created_at < ?", requestDto.CreatedBeforeAt.GetAsGoTime())
	}
	if requestDto.CreatedAfterAt != nil {
		dbQuery = dbQuery.Where("created_at > ?", requestDto.CreatedAfterAt.GetAsGoTime())
	}
	if requestDto.StartedBeforeAt != nil {
		dbQuery = dbQuery.Where("started_at < ?", requestDto.StartedBeforeAt.GetAsGoTime())
	}
	if requestDto.StartedAfterAt != nil {
		dbQuery = dbQuery.Where("started_at > ?", requestDto.StartedAfterAt.GetAsGoTime())
	}
	if requestDto.StoppedBeforeAt != nil {
		dbQuery = dbQuery.Where("stopped_at < ?", requestDto.StoppedBeforeAt.GetAsGoTime())
	}
	if requestDto.StoppedAfterAt != nil {
		dbQuery = dbQuery.Where("stopped_at > ?", requestDto.StoppedAfterAt.GetAsGoTime())
	}

	var itemsTotal int64
	err = dbQuery.Count(&itemsTotal).Error
	if err != nil {
		return responseDto, errors.New("CountItemsTotalError: " + err.Error())
	}

	dbQuery = dbQuery.Limit(int(requestDto.Pagination.ItemsPerPage))
	if requestDto.Pagination.LastSeenId == nil {
		offset := int(requestDto.Pagination.PageNumber) * int(requestDto.Pagination.ItemsPerPage)
		dbQuery = dbQuery.Offset(offset)
	} else {
		dbQuery = dbQuery.Where("id > ?", requestDto.Pagination.LastSeenId.String())
	}
	if requestDto.Pagination.SortBy != nil {
		orderStatement := requestDto.Pagination.SortBy.String()
		orderStatement = strcase.ToSnake(orderStatement)
		if orderStatement == "id" {
			orderStatement = "ID"
		}

		if requestDto.Pagination.SortDirection != nil {
			orderStatement += " " + requestDto.Pagination.SortDirection.String()
		}

		dbQuery = dbQuery.Order(orderStatement)
	}

	containerModels := []dbModel.Container{}
	err = dbQuery.Find(&containerModels).Error
	if err != nil {
		return responseDto, err
	}

	containerEntities := []entity.Container{}
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
		containerEntities = append(containerEntities, containerEntity)
	}

	itemsTotalUint := uint64(itemsTotal)
	pagesTotal := uint32(
		math.Ceil(float64(itemsTotal) / float64(requestDto.Pagination.ItemsPerPage)),
	)
	responsePagination := dto.Pagination{
		PageNumber:    requestDto.Pagination.PageNumber,
		ItemsPerPage:  requestDto.Pagination.ItemsPerPage,
		SortBy:        requestDto.Pagination.SortBy,
		SortDirection: requestDto.Pagination.SortDirection,
		PagesTotal:    &pagesTotal,
		ItemsTotal:    &itemsTotalUint,
	}

	responseDto = dto.ReadContainersResponse{
		Pagination:            responsePagination,
		Containers:            []entity.Container{},
		ContainersWithMetrics: []dto.ContainerWithMetrics{},
	}

	if requestDto.WithMetrics != nil && *requestDto.WithMetrics {
		containersWithMetrics, err := repo.containersWithMetricsFactory(containerEntities)
		if err != nil {
			return responseDto, err
		}

		responseDto.ContainersWithMetrics = containersWithMetrics
		return responseDto, nil
	}

	responseDto.Containers = containerEntities
	return responseDto, nil
}
