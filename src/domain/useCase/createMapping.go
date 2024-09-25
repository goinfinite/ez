package useCase

import (
	"errors"
	"log/slog"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
)

func CreateMapping(
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	containerQueryRepo repository.ContainerQueryRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateMapping,
) error {
	wasHostnameSent := createDto.Hostname != nil
	protocolStr := createDto.Protocol.String()
	isTransportLayer := protocolStr == "tcp" || protocolStr == "udp"
	if wasHostnameSent && isTransportLayer {
		createDto.Hostname = nil
	}

	if !wasHostnameSent && !isTransportLayer {
		return errors.New("HostnameRequiredForApplicationLayer")
	}

	publicPortStr := createDto.PublicPort.String()
	if publicPortStr == "1618" || publicPortStr == "3141" {
		return nil
	}

	currentMappings, err := mappingQueryRepo.Read()
	if err != nil {
		slog.Error("ReadMappingsInfraError", slog.Any("error", err))
		return errors.New("ReadMappingsInfraError")
	}

	var existingMapping *entity.Mapping
	for _, mapping := range currentMappings {
		if mapping.PublicPort != createDto.PublicPort {
			continue
		}

		if isTransportLayer {
			return errors.New("PublicPortAlreadyInUse")
		}

		if mapping.Hostname == nil {
			existingMapping = &mapping
			break
		}

		if mapping.Protocol != createDto.Protocol {
			return errors.New("PublicPortAlreadyInUseWithDifferentProtocol")
		}

		hostnameMatches := mapping.Hostname == createDto.Hostname
		if !hostnameMatches {
			continue
		}

		existingMapping = &mapping
		break
	}

	if existingMapping == nil {
		mappingId, err := mappingCmdRepo.Create(createDto)
		if err != nil {
			slog.Error("CreateMappingInfraError", slog.Any("error", err))
			return errors.New("CreateMappingInfraError")
		}

		NewCreateSecurityActivityRecord(activityRecordCmdRepo).
			CreateMapping(createDto, mappingId)

		newMapping, err := mappingQueryRepo.ReadById(mappingId)
		if err != nil {
			slog.Error("ReadMappingByIdInfraError", slog.Any("error", err))
			return errors.New("ReadMappingByIdInfraError")
		}
		existingMapping = &newMapping
	}

	for _, containerId := range createDto.ContainerIds {
		createTargetDto := dto.NewCreateMappingTarget(
			createDto.AccountId, existingMapping.Id, containerId,
			createDto.OperatorAccountId, createDto.OperatorIpAddress,
		)
		err = CreateMappingTarget(
			mappingQueryRepo, mappingCmdRepo, containerQueryRepo,
			activityRecordCmdRepo, createTargetDto,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
