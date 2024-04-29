package infra

import (
	"os"
	"testing"

	testHelpers "github.com/speedianet/control/src/devUtils"
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/valueObject"
)

func TestContainerCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	containerQueryRepo := NewContainerQueryRepo(persistentDbSvc)
	containerCmdRepo := NewContainerCmdRepo(persistentDbSvc)

	t.Run("CreateContainer", func(t *testing.T) {
		portBindings, _ := valueObject.NewPortBindingFromString("http")

		restartPolicy := valueObject.NewContainerRestartPolicyPanic("unless-stopped")

		profileId := valueObject.NewContainerProfileIdPanic(0)

		envs := []valueObject.ContainerEnv{
			valueObject.NewContainerEnvPanic("CONTROL_ENV1=testing"),
			valueObject.NewContainerEnvPanic("CONTROL_ENV2=testing"),
		}

		accountId := valueObject.NewAccountIdPanic(os.Getenv("DUMMY_USER_ID"))

		launchScript, err := valueObject.NewLaunchScript(
			"printenv > /tmp/hello.txt",
		)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}

		createContainer := dto.NewCreateContainer(
			accountId,
			valueObject.NewFqdnPanic("speedia.net"),
			valueObject.NewContainerImageAddressPanic("https://docker.io/nginx:latest"),
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
			t.Errorf("CreateContainerFailed: %v", err)
		}
	})

	t.Run("UpdateContainer", func(t *testing.T) {
		accId := valueObject.NewAccountIdPanic(os.Getenv("DUMMY_USER_ID"))
		containers, err := containerQueryRepo.GetByAccId(accId)
		if err != nil {
			t.Errorf("GetContainersFailed: %v", err)
		}

		if len(containers) == 0 {
			t.Error("NoContainersFound")
		}

		updateContainer := dto.NewUpdateContainer(
			accId,
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
		accId := valueObject.NewAccountIdPanic(os.Getenv("DUMMY_USER_ID"))
		containers, err := containerQueryRepo.GetByAccId(accId)
		if err != nil {
			t.Errorf("GetContainersFailed: %v", err)
		}

		if len(containers) == 0 {
			t.Error("NoContainersFound")
		}

		err = containerCmdRepo.Delete(
			accId,
			containers[0].Id,
		)
		if err != nil {
			t.Errorf("DeleteContainerFailed: %v", err)
		}
	})
}
