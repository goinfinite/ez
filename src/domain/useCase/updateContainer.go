package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

func UpdateContainer(
	containerQueryRepo repository.ContainerQueryRepo,
	containerCmdRepo repository.ContainerCmdRepo,
	accountQueryRepo repository.AccountQueryRepo,
	accountCmdRepo repository.AccountCmdRepo,
	containerProfileQueryRepo repository.ContainerProfileQueryRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	updateDto dto.UpdateContainer,
) error {
	readContainersDto := dto.ReadContainersRequest{
		Pagination:  ContainersDefaultPagination,
		ContainerId: []valueObject.ContainerId{updateDto.ContainerId},
	}

	responseDto, err := ReadContainers(containerQueryRepo, readContainersDto)
	if err != nil || len(responseDto.Containers) == 0 {
		return errors.New("ContainerNotFound")
	}
	containerEntity := responseDto.Containers[0]

	shouldUpdateQuota := updateDto.ProfileId != nil
	if shouldUpdateQuota {
		err = CheckAccountQuota(
			accountQueryRepo, containerProfileQueryRepo, updateDto.AccountId,
			*updateDto.ProfileId, &containerEntity.ProfileId,
		)
		if err != nil {
			return err
		}
	}

	err = containerCmdRepo.Update(updateDto)
	if err != nil {
		slog.Error("UpdateContainerInfraError", slog.Any("error", err))
		return errors.New("UpdateContainerInfraError")
	}

	if shouldUpdateQuota {
		err = accountCmdRepo.UpdateQuotaUsage(updateDto.AccountId)
		if err != nil {
			slog.Error("UpdateAccountQuotaInfraError", slog.Any("error", err))
			return errors.New("UpdateAccountQuotaError")
		}
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).UpdateContainer(updateDto)

	return nil
}
