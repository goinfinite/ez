package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/repository"
)

func UpdateResourceProfile(
	resourceProfileCmdRepo repository.ResourceProfileCmdRepo,
	updateResourceProfileDto dto.UpdateResourceProfile,
	containerQueryRepo repository.ContainerQueryRepo,
	containerCmdRepo repository.ContainerCmdRepo,
) error {
	err := resourceProfileCmdRepo.Update(updateResourceProfileDto)
	if err != nil {
		log.Printf("UpdateResourceProfileError: %s", err)
		return errors.New("UpdateResourceProfileInfraError")
	}

	shouldUpdateContainers := updateResourceProfileDto.BaseSpecs != nil
	if !shouldUpdateContainers {
		return nil
	}

	containers, err := containerQueryRepo.Get()
	if err != nil {
		log.Printf("GetContainersError: %s", err)
		return errors.New("GetContainersInfraError")
	}

	for _, container := range containers {
		if container.ResourceProfileId != updateResourceProfileDto.Id {
			continue
		}

		updateContainerDto := dto.NewUpdateContainer(
			container.AccountId,
			container.Id,
			container.Status,
			&updateResourceProfileDto.Id,
		)

		err := containerCmdRepo.Update(updateContainerDto)
		if err != nil {
			log.Printf("UpdateContainerAfterResourceProfileUpdateError: %s", err)
			continue
		}
	}

	return nil
}
