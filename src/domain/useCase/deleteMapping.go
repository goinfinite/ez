package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func DeleteMapping(
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	mappingId valueObject.MappingId,
) error {
	_, err := mappingQueryRepo.ReadById(mappingId)
	if err != nil {
		return errors.New("MappingNotFound")
	}

	err = mappingCmdRepo.Delete(mappingId)
	if err != nil {
		return errors.New("DeleteMappingError")
	}

	log.Printf("MappingId '%v' deleted.", mappingId)

	return nil
}
