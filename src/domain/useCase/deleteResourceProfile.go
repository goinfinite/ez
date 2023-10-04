package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/sfm/src/domain/repository"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

func DeleteResourceProfile(
	resourceProfileCmdRepo repository.ResourceProfileCmdRepo,
	profileId valueObject.ResourceProfileId,
	containerQueryRepo repository.ContainerQueryRepo,
	containerCmdRepo repository.ContainerCmdRepo,
) error {
	defaultResourceProfileId := valueObject.ResourceProfileId(0)
	if profileId == defaultResourceProfileId {
		return errors.New("CannotDeleteDefaultResourceProfile")
	}

	err := resourceProfileCmdRepo.Delete(profileId)
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
