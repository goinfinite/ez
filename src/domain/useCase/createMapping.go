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

	protocolStr := createDto.Protocol.String()
	isTransportLayer := protocolStr == "tcp" || protocolStr == "udp"

	if wasHostnameSent && isTransportLayer {
		createDto.Hostname = nil
	}

	publicPortStr := createDto.PublicPort.String()
	if publicPortStr == "1618" || publicPortStr == "3141" {
		return nil
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

		log.Printf("Mapping for port '%s/%s' added.", publicPortStr, protocolStr)
	}

	for _, containerId := range createDto.ContainerIds {
		addTargetDto := dto.NewCreateMappingTarget(mappingId, containerId)

		err = CreateMappingTarget(
			mappingQueryRepo,
			mappingCmdRepo,
			containerQueryRepo,
			addTargetDto,
		)
		if err != nil {
			log.Printf("[%s] CreateMappingTargetError: %s", containerId.String(), err)
			return errors.New("CreateMappingTargetInfraError")
		}
	}

	return nil
}
