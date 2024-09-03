package infra

import (
	"os"
	"testing"

	testHelpers "github.com/speedianet/control/src/devUtils"
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/valueObject"
)

var LocalOperatorAccountId, _ = valueObject.NewAccountId(0)
var LocalOperatorIpAddress = valueObject.NewLocalhostIpAddress()

func createDummyContainer(containerCmdRepo *ContainerCmdRepo) error {
	portBindings, _ := valueObject.NewPortBindingFromString("http")

	restartPolicy, _ := valueObject.NewContainerRestartPolicy("always")

	profileId, _ := valueObject.NewContainerProfileId(0)

	env1, _ := valueObject.NewContainerEnv("CONTROL_ENV1=testing")
	env2, _ := valueObject.NewContainerEnv("CONTROL_ENV2=testing")
	envs := []valueObject.ContainerEnv{env1, env2}

	accountId, _ := valueObject.NewAccountId(os.Getenv("DUMMY_USER_ID"))

	launchScript, err := valueObject.NewLaunchScript(
		"printenv > /tmp/hello.txt",
	)
	if err != nil {
		return err
	}

	containerHostname, _ := valueObject.NewFqdn("speedia.net")
	containerImage, _ := valueObject.NewContainerImageAddress("https://docker.io/speedianet/os")

	createContainer := dto.NewCreateContainer(
		accountId, containerHostname, containerImage, portBindings, &restartPolicy,
		nil, &profileId, envs, &launchScript, false,
		LocalOperatorAccountId, LocalOperatorIpAddress,
	)

	_, err = containerCmdRepo.Create(createContainer)
	return err
}

func deleteDummyContainer(
	containerQueryRepo *ContainerQueryRepo,
	containerCmdRepo *ContainerCmdRepo,
) error {
	containers, err := containerQueryRepo.Read()
	if err != nil {
		return err
	}

	if len(containers) == 0 {
		return nil
	}

	deleteDto := dto.NewDeleteContainer(
		containers[0].AccountId, containers[0].Id,
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
		containers, err := containerQueryRepo.ReadByAccountId(accountId)
		if err != nil {
			t.Errorf("ReadContainersFailed: %v", err)
		}

		if len(containers) == 0 {
			t.Error("NoContainersFound")
		}

		updateContainer := dto.NewUpdateContainer(
			accountId, containers[0].Id, nil, nil,
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
