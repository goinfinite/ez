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
	mappingId valueObject.MappingId,
	targetId valueObject.MappingTargetId,
) error {
	_, err := mappingQueryRepo.GetTargetById(targetId)
	if err != nil {
		return errors.New("MappingTargetNotFound")
	}

	err = mappingCmdRepo.DeleteTarget(mappingId, targetId)
	if err != nil {
		log.Printf("DeleteMappingTargetError: %v", err)
		return errors.New("DeleteMappingTargetInfraError")
	}

	log.Printf("TargetId '%v' deleted.", targetId)

	return nil
}
