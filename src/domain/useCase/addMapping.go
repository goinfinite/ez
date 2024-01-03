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
	addDto dto.AddMapping,
) error {
	if len(addDto.Targets) == 0 {
		return errors.New("NoTargetsToAdd")
	}

	wasHostnameSent := addDto.Hostname != nil

	isTcp := addDto.Protocol.String() == "tcp"
	isUdp := addDto.Protocol.String() == "udp"
	isTransportLayer := isTcp || isUdp

	if wasHostnameSent && isTransportLayer {
		return errors.New("TransportLayerCannotHaveHostname")
	}

	existingMapping, err := mappingQueryRepo.FindOne(
		addDto.Hostname,
		addDto.Port,
		addDto.Protocol,
	)
	if err != nil && err.Error() != "MappingNotFound" {
		log.Printf("FindExistingMappingError: %s", err)
		return errors.New("FindExistingMappingInfraError")
	}

	mappingExists := existingMapping == nil
	if !mappingExists {
		err = mappingCmdRepo.Add(addDto)
		if err != nil {
			log.Printf("AddMappingError: %s", err)
			return errors.New("AddMappingInfraError")
		}

		log.Printf(
			"Mapping for port '%v/%v' added.",
			addDto.Port,
			addDto.Protocol.String(),
		)
	}

	for _, target := range addDto.Targets {
		addTargetDto := dto.NewAddMappingTarget(
			existingMapping.Id,
			target.ContainerId,
			target.Port,
			target.Protocol,
		)

		err = AddMappingTarget(
			mappingQueryRepo,
			mappingCmdRepo,
			addTargetDto,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
