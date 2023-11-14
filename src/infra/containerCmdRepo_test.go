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
	dbSvc := testHelpers.GetDbSvc()
	containerQueryRepo := NewContainerQueryRepo(dbSvc)
	containerCmdRepo := NewContainerCmdRepo(dbSvc)

	t.Run("AddContainer", func(t *testing.T) {
		portBindings := []valueObject.PortBinding{
			valueObject.NewPortBinding(
				valueObject.NewNetworkProtocolPanic("tcp"),
				8080,
				8080,
			),
			valueObject.NewPortBinding(
				valueObject.NewNetworkProtocolPanic("tcp"),
				8443,
				8443,
			),
		}

		restartPolicy := valueObject.NewContainerRestartPolicyPanic("unless-stopped")

		profileId := valueObject.NewContainerProfileIdPanic(0)

		envs := []valueObject.ContainerEnv{
			valueObject.NewContainerEnvPanic("CONTROL_ENV1=testing"),
			valueObject.NewContainerEnvPanic("CONTROL_ENV2=testing"),
		}

		accountId := valueObject.NewAccountIdPanic(os.Getenv("DUMMY_USER_ID"))

		addContainer := dto.NewAddContainer(
			accountId,
			valueObject.NewFqdnPanic("speedia.net"),
			valueObject.NewContainerImgAddressPanic("https://docker.io/nginx:latest"),
			portBindings,
			&restartPolicy,
			nil,
			&profileId,
			envs,
		)

		err := containerCmdRepo.Add(addContainer)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
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

		err = containerCmdRepo.Update(containers[0], updateContainer)
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
