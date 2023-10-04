package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/repository"
)

func AddResourceProfile(
	resourceProfileCmdRepo repository.ResourceProfileCmdRepo,
	addResourceProfile dto.AddResourceProfile,
) error {
	err := resourceProfileCmdRepo.Add(addResourceProfile)
	if err != nil {
		log.Printf("AddResourceProfileError: %s", err)
		return errors.New("AddResourceProfileInfraError")
	}

	log.Printf("ResourceProfile '%v' added.", addResourceProfile.Name.String())

	return nil
}
