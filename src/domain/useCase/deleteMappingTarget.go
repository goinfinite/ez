package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func DeleteMappingTarget(
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	targetId valueObject.MappingTargetId,
) error {
	_, err := mappingQueryRepo.GetTargetById(targetId)
	if err != nil {
		return errors.New("MappingTargetNotFound")
	}

	err = mappingCmdRepo.DeleteTarget(targetId)
	if err != nil {
		return errors.New("DeleteMappingTargetError")
	}

	log.Printf("Mapping target '%v' deleted.", targetId)

	return nil
}
