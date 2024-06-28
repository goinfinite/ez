package infra

import (
	"testing"

	testHelpers "github.com/speedianet/control/src/devUtils"
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/valueObject"
)

func TestSecurityCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	securityCmdRepo := NewSecurityCmdRepo(persistentDbSvc)

	t.Run("CreateSecurityEvent", func(t *testing.T) {
		eventType, _ := valueObject.NewSecurityEventType("failed-login")
		ipAddress, _ := valueObject.NewIpAddress("127.0.0.1")
		createDto := dto.NewCreateSecurityEvent(eventType, nil, &ipAddress, nil)

		err := securityCmdRepo.CreateEvent(createDto)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})

	t.Run("DeleteSecurityEvents", func(t *testing.T) {
		ipAddress, _ := valueObject.NewIpAddress("127.0.0.1")
		deleteDto := dto.NewDeleteSecurityEvents(nil, &ipAddress, nil, nil)

		err := securityCmdRepo.DeleteEvents(deleteDto)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})
}
