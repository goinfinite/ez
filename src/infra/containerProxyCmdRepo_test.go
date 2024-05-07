package infra

import (
	"errors"
	"testing"

	testHelpers "github.com/speedianet/control/src/devUtils"
	"github.com/speedianet/control/src/domain/entity"
)

func getLastContainer(
	containerQueryRepo *ContainerQueryRepo,
) (containerEntity entity.Container, err error) {
	containers, err := containerQueryRepo.Get()
	if err != nil {
		return containerEntity, err
	}

	if len(containers) == 0 {
		return containerEntity, errors.New("NoContainersFound")
	}

	return containers[0], nil
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
		t.Errorf("GetContainersFailed: %v", err)
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
