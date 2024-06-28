package infra

import (
	"testing"

	testHelpers "github.com/speedianet/control/src/devUtils"
	"github.com/speedianet/control/src/domain/dto"
)

func TestSecurityQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	securityQueryRepo := NewSecurityQueryRepo(persistentDbSvc)

	t.Run("ReadSecurityEvents", func(t *testing.T) {
		readDto := dto.NewReadSecurityEvents(nil, nil, nil, nil)
		_, err := securityQueryRepo.ReadEvents(readDto)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})
}
