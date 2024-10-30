package infra

import (
	"os"
	"testing"

	testHelpers "github.com/goinfinite/ez/src/devUtils"
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/useCase"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

var LocalOperatorAccountId, _ = valueObject.NewAccountId(0)
var LocalOperatorIpAddress = valueObject.NewLocalhostIpAddress()

func createDummyContainer(containerCmdRepo *ContainerCmdRepo) error {
	portBindings, _ := valueObject.NewPortBindingFromString("http")

	restartPolicy, _ := valueObject.NewContainerRestartPolicy("always")

	profileId, _ := valueObject.NewContainerProfileId(0)

	env1, _ := valueObject.NewContainerEnv("Ez_ENV1=testing")
	env2, _ := valueObject.NewContainerEnv("Ez_ENV2=testing")
	envs := []valueObject.ContainerEnv{env1, env2}

	accountId, _ := valueObject.NewAccountId(os.Getenv("DUMMY_USER_ID"))

	launchScript, err := valueObject.NewLaunchScript(
		"printenv > /tmp/hello.txt",
	)
	if err != nil {
		return err
	}

	containerHostname, _ := valueObject.NewFqdn("goinfinite.net")
	imageAddress, _ := valueObject.NewContainerImageAddress("https://docker.io/goinfinite/os")

	createContainer := dto.NewCreateContainer(
		accountId, containerHostname, imageAddress, nil, portBindings, &restartPolicy,
		nil, &profileId, envs, &launchScript, false, false, nil,
		LocalOperatorAccountId, LocalOperatorIpAddress,
	)

	_, err = containerCmdRepo.Create(createContainer)
	return err
}

func deleteDummyContainer(
	containerQueryRepo *ContainerQueryRepo,
	containerCmdRepo *ContainerCmdRepo,
) error {
	requestDto := dto.ReadContainersRequest{
		Pagination: useCase.ContainersDefaultPagination,
	}

	responseDto, err := containerQueryRepo.Read(requestDto)
	if err != nil || len(responseDto.Containers) == 0 {
		return err
	}
	containerEntity := responseDto.Containers[0]

	deleteDto := dto.NewDeleteContainer(
		containerEntity.AccountId, containerEntity.Id,
		LocalOperatorAccountId, LocalOperatorIpAddress,
	)

	return containerCmdRepo.Delete(deleteDto)
}

func TestContainerCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	containerQueryRepo := NewContainerQueryRepo(persistentDbSvc)
	containerCmdRepo := NewContainerCmdRepo(persistentDbSvc)

	t.Run("CreateContainer", func(t *testing.T) {
		err := createDummyContainer(containerCmdRepo)
		if err != nil {
			t.Errorf("CreateContainerFailed: %v", err)
			return
		}
	})

	t.Run("UpdateContainer", func(t *testing.T) {
		accountId, _ := valueObject.NewAccountId(os.Getenv("DUMMY_USER_ID"))
		readContainersRequestDto := dto.ReadContainersRequest{
			Pagination:         useCase.ContainersDefaultPagination,
			ContainerAccountId: &accountId,
		}

		responseDto, err := containerQueryRepo.Read(readContainersRequestDto)
		if err != nil || len(responseDto.Containers) == 0 {
			t.Fatal(err)
		}
		containerEntity := responseDto.Containers[0]

		updateContainer := dto.NewUpdateContainer(
			containerEntity.AccountId, containerEntity.Id, nil, nil,
			LocalOperatorAccountId, LocalOperatorIpAddress,
		)

		err = containerCmdRepo.Update(updateContainer)
		if err != nil {
			t.Errorf("UpdateContainerFailed: %v", err)
		}
	})

	t.Run("DeleteContainer", func(t *testing.T) {
		err := deleteDummyContainer(containerQueryRepo, containerCmdRepo)
		if err != nil {
			t.Errorf("DeleteContainerFailed: %v", err)
		}
	})
}
