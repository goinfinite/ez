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
	addMapping dto.AddMapping,
) error {
	_, err := mappingQueryRepo.GetByHostPortProtocol(
		addMapping.Hostname,
		addMapping.Port,
		addMapping.Protocol,
	)
	if err == nil {
		return errors.New("MappingAlreadyExists")
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
