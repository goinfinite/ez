package useCase

import (
	"errors"
	"log"
	"slices"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func CreateMappingsWithContainerId(
	containerQueryRepo repository.ContainerQueryRepo,
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	containerId valueObject.ContainerId,
) error {
	containerEntity, err := containerQueryRepo.GetById(containerId)
	if err != nil {
		return errors.New("ContainerNotFound")
	}

	if containerEntity.IsSpeediaOs() {
		for portBindingIndex, portBinding := range containerEntity.PortBindings {
			if portBinding.PublicPort.String() != "1618" {
				continue
			}

			containerEntity.PortBindings = slices.Delete(
				containerEntity.PortBindings, portBindingIndex, portBindingIndex+1,
			)
			break
		}

		err = mappingCmdRepo.CreateContainerProxy(containerId)
		if err != nil {
			log.Printf("CreateContainerProxyMappingError: %s", err)
			return errors.New("CreateContainerProxyMappingInfraError")
		}
	}

	for _, portBinding := range containerEntity.PortBindings {
		createMappingDto := dto.NewCreateMapping(
			containerEntity.AccountId,
			&containerEntity.Hostname,
			portBinding.PublicPort,
			portBinding.Protocol,
			nil,
			nil,
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
