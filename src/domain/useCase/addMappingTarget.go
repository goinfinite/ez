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
	containerQueryRepo repository.ContainerQueryRepo,
	addDto dto.AddMappingTarget,
) error {
	mappingEntity, err := mappingQueryRepo.GetById(addDto.MappingId)
	if err != nil {
		log.Printf("GetMappingError: %s", err)
		return errors.New("GetMappingInfraError")
	}

	containerEntity, err := containerQueryRepo.GetById(addDto.ContainerId)
	if err != nil {
		log.Printf("GetContainerError: %s", err)
		return errors.New("GetContainerInfraError")
	}

	publicPortMatches := false
	for _, portBinding := range containerEntity.PortBindings {
		if portBinding.PublicPort != mappingEntity.PublicPort {
			continue
		}
		publicPortMatches = true
	}

	if !publicPortMatches {
		log.Printf(
			"ContainerId '%s' does not bind to public port '%d'.",
			addDto.ContainerId, mappingEntity.PublicPort,
		)
		return errors.New("ContainerDoesNotBindToMappingPublicPort")
	}

	for _, target := range mappingEntity.Targets {
		if target.ContainerId != addDto.ContainerId {
			continue
		}

		log.Printf(
			"Skipping AddMappingTarget: ContainerId '%s' is already a target for MappingId '%s'.",
			addDto.ContainerId, addDto.MappingId,
		)
		return nil
	}

	err = mappingCmdRepo.AddTarget(addDto)
	if err != nil {
		log.Printf("AddMappingTargetError: %s", err)
		return errors.New("AddMappingTargetInfraError")
	}

	log.Printf(
		"ContainerId '%s' added as target for MappingId '%s'.",
		addDto.ContainerId, addDto.MappingId,
	)

	return nil
}
