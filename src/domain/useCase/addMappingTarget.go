package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
)

func AddMappingTarget(
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	addDto dto.AddMappingTarget,
) error {
	_, err := mappingQueryRepo.GetById(addDto.MappingId)
	if err != nil {
		log.Printf("GetMappingError: %s", err)
		return errors.New("GetMappingInfraError")
	}

	err = mappingCmdRepo.AddTarget(addDto)
	if err != nil {
		log.Printf("AddMappingTargetError: %s", err)
		return errors.New("AddMappingTargetInfraError")
	}

	log.Printf(
		"'%s' added as target for mapping with ID '%s'.",
		addDto.ContainerId,
		addDto.MappingId,
	)

	return nil
}
