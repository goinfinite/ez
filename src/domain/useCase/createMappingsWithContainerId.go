package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func CreateMappingsWithContainerId(
	containerQueryRepo repository.ContainerQueryRepo,
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	containerProxyCmdRepo repository.ContainerProxyCmdRepo,
	containerId valueObject.ContainerId,
) error {
	containerEntity, err := containerQueryRepo.GetById(containerId)
	if err != nil {
		return errors.New("ContainerNotFound")
	}

	for _, portBinding := range containerEntity.PortBindings {
		publicPortStr := portBinding.PublicPort.String()
		if publicPortStr == "1618" || publicPortStr == "3141" {
			continue
		}

		createMappingDto := dto.NewCreateMapping(
			containerEntity.AccountId,
			&containerEntity.Hostname,
			portBinding.PublicPort,
			portBinding.Protocol,
			[]valueObject.ContainerId{containerId},
		)
		err = CreateMapping(
			mappingQueryRepo,
			mappingCmdRepo,
			containerQueryRepo,
			createMappingDto,
		)
		if err != nil {
			log.Printf("CreateMappingError: %s", err)
			return errors.New("CreateMappingInfraError")
		}
	}

	return nil
}
