package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
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

	currentMappings, err := mappingQueryRepo.Get()
	if err != nil {
		log.Printf("GetMappingError: %s", err)
		return errors.New("GetMappingInfraError")
	}

	var existingMapping *entity.Mapping
	for _, mapping := range currentMappings {
		if mapping.PublicPort != createDto.PublicPort {
			continue
		}

		if isTransportLayer {
			return errors.New("PublicPortAlreadyInUse")
		}

		hostnameMatches := mapping.Hostname == createDto.Hostname
		if hostnameMatches && mapping.Protocol != createDto.Protocol {
			return errors.New("PublicPortAlreadyInUseWithDifferentProtocol")
		}

		existingMapping = &mapping
		break
	}

	if existingMapping == nil {
		mappingId, err := mappingCmdRepo.Create(createDto)
		if err != nil {
			log.Printf("CreateMappingError: %s", err)
			return errors.New("CreateMappingInfraError")
		}

		log.Printf("Mapping for port '%s/%s' added.", publicPortStr, protocolStr)

		newMapping, err := mappingQueryRepo.GetById(mappingId)
		if err != nil {
			log.Printf("GetMappingByIdError: %s", err)
			return errors.New("GetMappingByIdInfraError")
		}
		existingMapping = &newMapping
	}

	for _, containerId := range createDto.ContainerIds {
		createTargetDto := dto.NewCreateMappingTarget(existingMapping.Id, containerId)
		err = CreateMappingTarget(
			mappingQueryRepo,
			mappingCmdRepo,
			containerQueryRepo,
			createTargetDto,
		)
		if err != nil {
			log.Printf("[%s] CreateMappingTargetError: %s", containerId.String(), err)
			return errors.New("CreateMappingTargetInfraError")
		}
	}

	return nil
}
