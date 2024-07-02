package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
)

func StopAllContainers(
	containerQueryRepo repository.ContainerQueryRepo,
	containerCmdRepo repository.ContainerCmdRepo,
) error {
	containers, err := containerQueryRepo.Read()
	if err != nil {
		return errors.New("ReadContainersError: " + err.Error())
	}

	for _, currentContainer := range containers {
		newContainerStatus := false
		updateDto := dto.NewUpdateContainer(
			currentContainer.AccountId,
			currentContainer.Id,
			&newContainerStatus,
			&currentContainer.ProfileId,
		)

		err = containerCmdRepo.Update(updateDto)
		if err != nil {
			log.Printf(
				"[ContainerId: %s] StopContainerError: %s",
				currentContainer.Id,
				err,
			)
			continue
		}
	}

	return nil
}
