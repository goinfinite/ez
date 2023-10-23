package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/sfm/src/domain/entity"
	"github.com/speedianet/sfm/src/domain/repository"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

func DeleteContainerProfile(
	containerProfileQueryRepo repository.ContainerProfileQueryRepo,
	containerProfileCmdRepo repository.ContainerProfileCmdRepo,
	containerQueryRepo repository.ContainerQueryRepo,
	containerCmdRepo repository.ContainerCmdRepo,
	profileId valueObject.ContainerProfileId,
) error {
	_, err := containerProfileQueryRepo.GetById(profileId)
	if err != nil {
		return errors.New("ContainerProfileNotFound")
	}

	defaultContainerProfileId := entity.DefaultContainerProfile().Id
	if profileId == defaultContainerProfileId {
		return errors.New("CannotDeleteDefaultContainerProfile")
	}

	err = containerProfileCmdRepo.Delete(profileId)
	if err != nil {
		log.Printf("DeleteContainerProfileError: %s", err)
		return errors.New("DeleteContainerProfileInfraError")
	}

	err = updateContainerContainerProfileId(
		containerQueryRepo,
		containerCmdRepo,
		defaultContainerProfileId,
	)
	if err != nil {
		log.Printf("UpdateContainerContainerProfileError: %s", err)
		return errors.New("UpdateContainerContainerProfileInfraError")
	}

	return nil
}
