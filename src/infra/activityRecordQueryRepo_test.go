package infra

import (
	"testing"

	testHelpers "github.com/speedianet/control/src/devUtils"
	"github.com/speedianet/control/src/domain/dto"
)

func TestActivityRecordQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	trailDbSvc := testHelpers.GetTrailDbSvc()
	activityRecordQueryRepo := NewActivityRecordQueryRepo(trailDbSvc)

	t.Run("ReadActivityRecordQuery", func(t *testing.T) {
		readDto := dto.ReadActivityRecords{}
		_, err := activityRecordQueryRepo.Read(readDto)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})
}
