package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func DeleteContainerProfile(
	containerProfileQueryRepo repository.ContainerProfileQueryRepo,
	containerProfileCmdRepo repository.ContainerProfileCmdRepo,
	containerQueryRepo repository.ContainerQueryRepo,
	containerCmdRepo repository.ContainerCmdRepo,
	profileId valueObject.ContainerProfileId,
) error {
	_, err := containerProfileQueryRepo.ReadById(profileId)
	if err != nil {
		return errors.New("ContainerProfileNotFound")
	}

	defaultContainerProfileId := entity.DefaultContainerProfile().Id
	if profileId == defaultContainerProfileId {
		return errors.New("CannotDeleteDefaultContainerProfile")
	}

	err = updateContainersWithProfileId(
		containerQueryRepo,
		containerCmdRepo,
		defaultContainerProfileId,
	)
	if err != nil {
		log.Printf("UpdateContainersBackToDefaultProfileError: %s", err)
		return errors.New("UpdateContainersBackToDefaultProfileInfraError")
	}

	err = containerProfileCmdRepo.Delete(profileId)
	if err != nil {
		log.Printf("DeleteContainerProfileError: %s", err)
		return errors.New("DeleteContainerProfileInfraError")
	}

	return nil
}
