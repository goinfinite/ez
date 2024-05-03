package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
)

func CreateMappingTarget(
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	containerQueryRepo repository.ContainerQueryRepo,
	createDto dto.CreateMappingTarget,
) error {
	mappingEntity, err := mappingQueryRepo.GetById(createDto.MappingId)
	if err != nil {
		log.Printf("GetMappingError: %s", err)
		return errors.New("GetMappingInfraError")
	}

	containerEntity, err := containerQueryRepo.GetById(createDto.ContainerId)
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
			createDto.ContainerId, mappingEntity.PublicPort,
		)
		return errors.New("ContainerDoesNotBindToMappingPublicPort")
	}

	for _, target := range mappingEntity.Targets {
		if target.ContainerId != createDto.ContainerId {
			continue
		}

		log.Printf(
			"Skipping CreateMappingTarget: ContainerId '%s' is already a target for MappingId '%s'.",
			createDto.ContainerId, createDto.MappingId,
		)
		return nil
	}

	err = mappingCmdRepo.CreateTarget(createDto)
	if err != nil {
		log.Printf("CreateMappingTargetError: %s", err)
		return errors.New("CreateMappingTargetInfraError")
	}

	log.Printf(
		"ContainerId '%s' created as target for MappingId '%s'.",
		createDto.ContainerId, createDto.MappingId,
	)

	return nil
}
