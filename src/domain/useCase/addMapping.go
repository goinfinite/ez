package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
)

func AddMapping(
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	containerQueryRepo repository.ContainerQueryRepo,
	addDto dto.AddMapping,
) error {
	wasHostnameSent := addDto.Hostname != nil

	isTcp := addDto.Protocol.String() == "tcp"
	isUdp := addDto.Protocol.String() == "udp"
	isTransportLayer := isTcp || isUdp

	if wasHostnameSent && isTransportLayer {
		addDto.Hostname = nil
	}

	existingMapping, err := mappingQueryRepo.FindOne(
		addDto.Hostname,
		addDto.PublicPort,
		addDto.Protocol,
	)
	if err != nil && err.Error() != "MappingNotFound" {
		log.Printf("FindExistingMappingError: %s", err)
		return errors.New("FindExistingMappingInfraError")
	}

	mappingId := existingMapping.Id
	mappingAlreadyExists := mappingId != 0
	if !mappingAlreadyExists {
		mappingId, err = mappingCmdRepo.Add(addDto)
		if err != nil {
			log.Printf("AddMappingError: %s", err)
			return errors.New("AddMappingInfraError")
		}

		log.Printf(
			"Mapping for port '%v/%v' added.",
			addDto.PublicPort,
			addDto.Protocol.String(),
		)
	}

	for _, containerId := range addDto.ContainerIds {
		addTargetDto := dto.NewAddMappingTarget(
			mappingId,
			containerId,
		)

		err = AddMappingTarget(
			mappingQueryRepo,
			mappingCmdRepo,
			containerQueryRepo,
			addTargetDto,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
