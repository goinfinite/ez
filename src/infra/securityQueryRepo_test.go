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

	t.Run("GetSecurityEvents", func(t *testing.T) {
		getDto := dto.NewGetSecurityEvents(nil, nil, nil, nil)
		_, err := securityQueryRepo.GetEvents(getDto)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})
}
