package infra

import (
	"testing"

	testHelpers "github.com/goinfinite/ez/src/devUtils"
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/useCase"
)

func getLastContainer(
	containerQueryRepo *ContainerQueryRepo,
) (containerEntity entity.Container, err error) {
	readContainersRequestDto := dto.ReadContainersRequest{
		Pagination: useCase.ContainersDefaultPagination,
	}

	readContainersResponseDto, err := containerQueryRepo.Read(readContainersRequestDto)
	if err != nil || len(readContainersResponseDto.Containers) == 0 {
		return containerEntity, err
	}
	return readContainersResponseDto.Containers[0], nil
}

func TestContainerProxyCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	containerProxyCmdRepo := NewContainerProxyCmdRepo(persistentDbSvc)

	containerQueryRepo := NewContainerQueryRepo(persistentDbSvc)
	containerCmdRepo := NewContainerCmdRepo(persistentDbSvc)
	err := createDummyContainer(containerCmdRepo)
	if err != nil {
		t.Errorf("CreateContainerFailed: %v", err)
		return
	}

	lastContainer, err := getLastContainer(containerQueryRepo)
	if err != nil {
		t.Errorf("ReadContainersFailed: %v", err)
		return
	}

	t.Run("CreateContainerProxy", func(t *testing.T) {
		err = containerProxyCmdRepo.Create(lastContainer.Id)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
			return
		}
	})

	t.Run("DeleteContainerProxy", func(t *testing.T) {
		err = containerProxyCmdRepo.Delete(lastContainer.Id)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
			return
		}
	})

	err = deleteDummyContainer(containerQueryRepo, containerCmdRepo)
	if err != nil {
		t.Errorf("DeleteContainerFailed: %v", err)
		return
	}
}
