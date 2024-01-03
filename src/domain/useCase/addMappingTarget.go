package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func AddMappingTarget(
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	addDto dto.AddMappingTarget,
) error {
	mapping, err := mappingQueryRepo.GetById(addDto.MappingId)
	if err != nil {
		log.Printf("GetMappingError: %s", err)
		return errors.New("GetMappingInfraError")
	}

	err = addNewMappingTargets(
		mappingCmdRepo,
		&mapping,
		[]valueObject.MappingTarget{addDto.Target},
	)
	if err != nil {
		log.Printf("AddMappingTargetsError: %s", err)
		return errors.New("AddMappingTargetsInfraError")
	}

	log.Printf(
		"Target '%s/%s/%s' added for mapping with ID '%s'.",
		addDto.Target.ContainerId,
		addDto.Target.Port.String(),
		addDto.Target.Protocol.String(),
		addDto.MappingId,
	)

	return nil
}
