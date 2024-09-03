package useCase

import (
	"log/slog"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

var LocalOperatorAccountId, _ = valueObject.NewAccountId(0)
var LocalOperatorIpAddress = valueObject.NewLocalhostIpAddress()

func BootContainers(
	containerQueryRepo repository.ContainerQueryRepo,
	containerCmdRepo repository.ContainerCmdRepo,
) {
	containers, err := containerQueryRepo.Read()
	if err != nil {
		slog.Error("ReadContainersInfraError", slog.Any("error", err))
		return
	}

	for _, currentContainer := range containers {
		restartPolicyStr := currentContainer.RestartPolicy.String()
		shouldBoot := restartPolicyStr == "always" || restartPolicyStr == "on-failure"
		if !shouldBoot {
			continue
		}

		newContainerStatus := true
		updateDto := dto.NewUpdateContainer(
			currentContainer.AccountId, currentContainer.Id, &newContainerStatus,
			nil, LocalOperatorAccountId, LocalOperatorIpAddress,
		)

		err = containerCmdRepo.Update(updateDto)
		if err != nil {
			slog.Error(
				"UpdateContainerInfraError",
				slog.String("containerId", currentContainer.Id.String()),
				slog.Any("error", err),
			)
			continue
		}
	}
}
