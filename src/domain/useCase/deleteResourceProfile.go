package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/sfm/src/domain/repository"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

func DeleteResourceProfile(
	resourceProfileQueryRepo repository.ResourceProfileQueryRepo,
	resourceProfileCmdRepo repository.ResourceProfileCmdRepo,
	containerQueryRepo repository.ContainerQueryRepo,
	containerCmdRepo repository.ContainerCmdRepo,
	profileId valueObject.ResourceProfileId,
) error {
	_, err := resourceProfileQueryRepo.GetById(profileId)
	if err != nil {
		return errors.New("ResourceProfileNotFound")
	}

	defaultResourceProfileId := valueObject.ResourceProfileId(1)
	if profileId == defaultResourceProfileId {
		return errors.New("CannotDeleteDefaultResourceProfile")
	}

	err = resourceProfileCmdRepo.Delete(profileId)
	if err != nil {
		log.Printf("DeleteResourceProfileError: %s", err)
		return errors.New("DeleteResourceProfileInfraError")
	}

	err = updateContainerResourceProfileId(
		containerQueryRepo,
		containerCmdRepo,
		defaultResourceProfileId,
	)
	if err != nil {
		log.Printf("UpdateContainerResourceProfileError: %s", err)
		return errors.New("UpdateContainerResourceProfileInfraError")
	}

	return nil
}
