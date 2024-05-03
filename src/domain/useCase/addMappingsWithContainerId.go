package useCase

import (
	"errors"
	"log"
	"slices"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func AddMappingsWithContainerId(
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
			log.Printf("AddContainerProxyMappingError: %s", err)
			return errors.New("AddContainerProxyMappingInfraError")
		}
	}

	for _, portBinding := range containerEntity.PortBindings {
		addMappingDto := dto.NewAddMapping(
			containerEntity.AccountId,
			&containerEntity.Hostname,
			portBinding.PublicPort,
			portBinding.Protocol,
			[]valueObject.ContainerId{containerId},
		)
		err = AddMapping(
			mappingQueryRepo,
			mappingCmdRepo,
			containerQueryRepo,
			addMappingDto,
		)
		if err != nil {
			log.Printf("AddMappingError: %s", err)
			return errors.New("AddMappingInfraError")
		}
	}

	return nil
}
