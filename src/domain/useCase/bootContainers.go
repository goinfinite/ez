package useCase

import (
	"log"

	"github.com/goinfinite/fleet/src/domain/dto"
	"github.com/goinfinite/fleet/src/domain/repository"
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
		shouldBoot := restartPolicyStr == "always" || restartPolicyStr == "unless-stopped"
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

		err = containerCmdRepo.Update(currentContainer, updateDto)
		if err != nil {
			log.Printf(
				"[ContainerId: %s] UpdateContainerError: %s",
				currentContainer.Id,
				err,
			)
			continue
		}
	}
}
