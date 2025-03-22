package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

func DeleteContainer(
	containerQueryRepo repository.ContainerQueryRepo,
	containerCmdRepo repository.ContainerCmdRepo,
	accountCmdRepo repository.AccountCmdRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	containerProxyCmdRepo repository.ContainerProxyCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteContainer,
) error {
	readContainersDto := dto.ReadContainersRequest{
		Pagination:  ContainersDefaultPagination,
		ContainerId: []valueObject.ContainerId{deleteDto.ContainerId},
	}

	responseDto, err := ReadContainers(containerQueryRepo, readContainersDto)
	if err != nil || len(responseDto.Containers) == 0 {
		return errors.New("ContainerNotFound")
	}

	err = mappingCmdRepo.DeleteTargetsByContainerId(deleteDto.ContainerId)
	if err != nil {
		slog.Error("DeleteMappingTargetInfraError", slog.Any("error", err))
		return errors.New("DeleteMappingTargetInfraError")
	}

	err = mappingCmdRepo.DeleteEmpty()
	if err != nil {
		slog.Error("DeleteEmptyMappingInfraError", slog.Any("error", err))
		return errors.New("DeleteEmptyMappingInfraError")
	}

	err = containerProxyCmdRepo.Delete(deleteDto.ContainerId)
	if err != nil {
		slog.Error("DeleteContainerProxyInfraError", slog.Any("error", err))
		return errors.New("DeleteContainerProxyInfraError")
	}

	err = containerCmdRepo.Delete(deleteDto)
	if err != nil {
		slog.Error("DeleteContainerInfraError", slog.Any("error", err))
		return errors.New("DeleteContainerInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).DeleteContainer(deleteDto)

	err = accountCmdRepo.UpdateQuotaUsage(deleteDto.AccountId)
	if err != nil {
		slog.Error("UpdateAccountQuotaInfraError", slog.Any("error", err))
		return errors.New("UpdateAccountQuotaInfraError")
	}

	return nil
}
