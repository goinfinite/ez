package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/repository"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

func updateContainerResourceProfileId(
	containerQueryRepo repository.ContainerQueryRepo,
	containerCmdRepo repository.ContainerCmdRepo,
	profileId valueObject.ResourceProfileId,
) error {
	containers, err := containerQueryRepo.Get()
	if err != nil {
		log.Printf("GetContainersError: %s", err)
		return errors.New("GetContainersInfraError")
	}

	for _, container := range containers {
		if container.ResourceProfileId != profileId {
			continue
		}

		updateContainerDto := dto.NewUpdateContainer(
			container.AccountId,
			container.Id,
			container.Status,
			&profileId,
		)

		err := containerCmdRepo.Update(updateContainerDto)
		if err != nil {
			log.Printf("UpdateContainerResourceProfileError: %s", err)
			continue
		}
	}

	return nil
}

func UpdateResourceProfile(
	resourceProfileQueryRepo repository.ResourceProfileQueryRepo,
	resourceProfileCmdRepo repository.ResourceProfileCmdRepo,
	containerQueryRepo repository.ContainerQueryRepo,
	containerCmdRepo repository.ContainerCmdRepo,
	updateResourceProfileDto dto.UpdateResourceProfile,
) error {
	_, err := resourceProfileQueryRepo.GetById(updateResourceProfileDto.Id)
	if err != nil {
		return errors.New("ResourceProfileNotFound")
	}

	err = resourceProfileCmdRepo.Update(updateResourceProfileDto)
	if err != nil {
		log.Printf("UpdateResourceProfileError: %s", err)
		return errors.New("UpdateResourceProfileInfraError")
	}

	shouldUpdateContainers := updateResourceProfileDto.BaseSpecs != nil
	if !shouldUpdateContainers {
		return nil
	}

	err = updateContainerResourceProfileId(
		containerQueryRepo,
		containerCmdRepo,
		updateResourceProfileDto.Id,
	)
	if err != nil {
		log.Printf("UpdateResourceProfileContainersError: %s", err)
		return errors.New("UpdateResourceProfileContainersInfraError")
	}

	return nil
}
