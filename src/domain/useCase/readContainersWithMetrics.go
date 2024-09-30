package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
)

func ReadContainersWithMetrics(
	containerQueryRepo repository.ContainerQueryRepo,
) ([]dto.ContainerWithMetrics, error) {
	containersWithMetricsList, err := containerQueryRepo.ReadWithMetrics()
	if err != nil {
		slog.Error("ReadContainersWithMetricsInfraError", slog.Any("error", err))
		return containersWithMetricsList, errors.New("ReadContainersWithMetricsInfraError")
	}

	return containersWithMetricsList, nil
}
