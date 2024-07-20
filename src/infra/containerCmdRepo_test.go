package infra

import (
	"os"
	"testing"

	testHelpers "github.com/speedianet/control/src/devUtils"
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/valueObject"
)

func createDummyContainer(containerCmdRepo *ContainerCmdRepo) error {
	portBindings, _ := valueObject.NewPortBindingFromString("http")

	restartPolicy := valueObject.NewContainerRestartPolicyPanic("unless-stopped")

	profileId := valueObject.NewContainerProfileIdPanic(0)

	envs := []valueObject.ContainerEnv{
		valueObject.NewContainerEnvPanic("CONTROL_ENV1=testing"),
		valueObject.NewContainerEnvPanic("CONTROL_ENV2=testing"),
	}

	accountId, _ := valueObject.NewAccountId(os.Getenv("DUMMY_USER_ID"))

	launchScript, err := valueObject.NewLaunchScript(
		"printenv > /tmp/hello.txt",
	)
	if err != nil {
		return err
	}

	createContainer := dto.NewCreateContainer(
		accountId,
		valueObject.NewFqdnPanic("speedia.net"),
		valueObject.NewContainerImageAddressPanic("https://docker.io/speedianet/os"),
		portBindings,
		&restartPolicy,
		nil,
		&profileId,
		envs,
		&launchScript,
		false,
	)

	_, err = containerCmdRepo.Create(createContainer)
	if err != nil {
		return err
	}

	return nil
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

	return containerCmdRepo.Delete(
		containers[0].AccountId,
		containers[0].Id,
	)
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
			accountId,
			containers[0].Id,
			nil,
			nil,
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
