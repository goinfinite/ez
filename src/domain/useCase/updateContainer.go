package useCase

import (
	"errors"

	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/repository"
)

func UpdateContainer(
	containerQueryRepo repository.ContainerQueryRepo,
	containerCmdRepo repository.ContainerCmdRepo,
	accQueryRepo repository.AccQueryRepo,
	updateContainerDto dto.UpdateContainer,
) error {
	_, err := containerQueryRepo.GetById(
		updateContainerDto.AccountId,
		updateContainerDto.ContainerId,
	)
	if err != nil {
		return errors.New("ContainerNotFound")
	}

	err = AccQuotaCheck(
		accQueryRepo,
		updateContainerDto.AccountId,
		*updateContainerDto.BaseSpecs,
	)
	if err != nil {
		return err
	}

	return containerCmdRepo.Update(updateContainerDto)
}
