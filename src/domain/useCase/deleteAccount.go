package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

func DeleteAccount(
	accountQueryRepo repository.AccountQueryRepo,
	accountCmdRepo repository.AccountCmdRepo,
	containerQueryRepo repository.ContainerQueryRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteAccount,
) error {
	_, err := accountQueryRepo.ReadFirst(dto.ReadAccountsRequest{
		AccountId: &deleteDto.AccountId,
	})
	if err != nil {
		return errors.New("AccountNotFound")
	}

	readContainersDto := dto.ReadContainersRequest{
		Pagination:         ContainersDefaultPagination,
		ContainerAccountId: []valueObject.AccountId{deleteDto.AccountId},
	}

	responseDto, err := ReadContainers(containerQueryRepo, readContainersDto)
	if err != nil {
		return errors.New("ReadContainersInfraError")
	}

	if len(responseDto.Containers) > 0 {
		return errors.New("AccountStillHasContainers")
	}

	err = accountCmdRepo.Delete(deleteDto.AccountId)
	if err != nil {
		slog.Error("DeleteAccountInfraError", slog.Any("error", err))
		return errors.New("DeleteAccountInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).DeleteAccount(deleteDto)

	return nil
}
