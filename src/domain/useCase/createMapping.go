package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
)

func CreateMapping(
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	containerQueryRepo repository.ContainerQueryRepo,
	createDto dto.CreateMapping,
) error {
	wasHostnameSent := createDto.Hostname != nil

	isTcp := createDto.Protocol.String() == "tcp"
	isUdp := createDto.Protocol.String() == "udp"
	isTransportLayer := isTcp || isUdp

	if wasHostnameSent && isTransportLayer {
		createDto.Hostname = nil
	}

	existingMapping, err := mappingQueryRepo.FindOne(
		createDto.Hostname,
		createDto.PublicPort,
		createDto.Protocol,
	)
	if err != nil && err.Error() != "MappingNotFound" {
		log.Printf("FindExistingMappingError: %s", err)
		return errors.New("FindExistingMappingInfraError")
	}

	mappingId := existingMapping.Id
	mappingAlreadyExists := mappingId != 0
	if !mappingAlreadyExists {
		mappingId, err = mappingCmdRepo.Create(createDto)
		if err != nil {
			log.Printf("CreateMappingError: %s", err)
			return errors.New("CreateMappingInfraError")
		}

		log.Printf(
			"Mapping for port '%v/%v' added.",
			createDto.PublicPort,
			createDto.Protocol.String(),
		)
	}

	for _, containerId := range createDto.ContainerIds {
		addTargetDto := dto.NewCreateMappingTarget(
			mappingId,
			containerId,
		)

		err = CreateMappingTarget(
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