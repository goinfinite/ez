package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func mappingTargetAlreadyExists(
	mapping *entity.Mapping,
	target valueObject.MappingTarget,
) bool {
	if len(mapping.Targets) == 0 {
		return false
	}

	for _, currentTarget := range mapping.Targets {
		isSameContainerId := currentTarget.ContainerId == target.ContainerId
		isSamePort := currentTarget.Port == target.Port
		if isSameContainerId && isSamePort {
			return true
		}
	}

	return false
}

func addNewMappingTargets(
	mappingCmdRepo repository.MappingCmdRepo,
	existingMapping *entity.Mapping,
	newTargets []valueObject.MappingTarget,
) error {
	targetsToAdd := []valueObject.MappingTarget{}
	for _, newTarget := range newTargets {
		if mappingTargetAlreadyExists(existingMapping, newTarget) {
			continue
		}

		targetsToAdd = append(targetsToAdd, newTarget)
	}

	if len(targetsToAdd) == 0 {
		return errors.New("NoNewTargetsToAdd")
	}

	err := mappingCmdRepo.AddTargets(
		existingMapping.Id,
		targetsToAdd,
	)
	if err != nil {
		log.Printf("AddMappingTargetsError: %s", err)
		return errors.New("AddMappingTargetsInfraError")
	}

	return nil
}

func AddMapping(
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	addMapping dto.AddMapping,
) error {
	if len(addMapping.Targets) == 0 {
		return errors.New("NoTargetsToAdd")
	}

	wasHostnameSent := addMapping.Hostname != nil

	isTcp := addMapping.Protocol.String() == "tcp"
	isUdp := addMapping.Protocol.String() == "udp"
	isTransportLayer := isTcp || isUdp

	if wasHostnameSent && isTransportLayer {
		return errors.New("TransportLayerCannotHaveHostname")
	}

	existingMapping, err := mappingQueryRepo.FindOne(
		addMapping.Hostname,
		addMapping.Port,
		addMapping.Protocol,
	)
	if err != nil && err.Error() != "MappingNotFound" {
		log.Printf("FindExistingMappingError: %s", err)
		return errors.New("FindExistingMappingInfraError")
	}

	if existingMapping != nil {
		return addNewMappingTargets(
			mappingCmdRepo,
			existingMapping,
			addMapping.Targets,
		)
	}

	err = mappingCmdRepo.Add(addMapping)
	if err != nil {
		log.Printf("AddMappingError: %s", err)
		return errors.New("AddMappingInfraError")
	}

	log.Printf(
		"Mapping for port '%v/%v' added.",
		addMapping.Port,
		addMapping.Protocol.String(),
	)

	return nil
}
