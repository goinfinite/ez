package useCase

import (
	"errors"
	"log"

	"github.com/goinfinite/fleet/src/domain/dto"
	"github.com/goinfinite/fleet/src/domain/repository"
)

func AddContainerProfile(
	containerProfileCmdRepo repository.ContainerProfileCmdRepo,
	addContainerProfile dto.AddContainerProfile,
) error {
	err := containerProfileCmdRepo.Add(addContainerProfile)
	if err != nil {
		log.Printf("AddContainerProfileError: %s", err)
		return errors.New("AddContainerProfileInfraError")
	}

	log.Printf("ContainerProfile '%v' added.", addContainerProfile.Name.String())

	return nil
}
