package useCase

import (
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
)

func BootContainers(
	containerQueryRepo repository.ContainerQueryRepo,
	containerCmdRepo repository.ContainerCmdRepo,
) {
	containers, err := containerQueryRepo.Get()
	if err != nil {
		log.Printf("GetContainersError: %v", err)
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
			currentContainer.AccountId,
			currentContainer.Id,
			&newContainerStatus,
			&currentContainer.ProfileId,
		)

		err = containerCmdRepo.Update(updateDto)
		if err != nil {
			log.Printf(
				"[ContainerId: %s] UpdateContainerError: %s", currentContainer.Id, err,
			)
			continue
		}
	}
}
