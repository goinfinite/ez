package presenterHelper

import (
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/infra/db"
	"github.com/goinfinite/ez/src/presentation/service"
	componentForm "github.com/goinfinite/ez/src/presentation/ui/component/form"
)

func ReadBackupDestinationSelectLabelValuePairs(
	persistentDbSvc *db.PersistentDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) []componentForm.SelectLabelValuePair {
	selectLabelValuePairs := []componentForm.SelectLabelValuePair{}

	backupService := service.NewBackupService(persistentDbSvc, trailDbSvc)

	readBackupDestinationsServiceOutput := backupService.ReadDestination(map[string]interface{}{
		"itemsPerPage": 1000,
	})
	if readBackupDestinationsServiceOutput.Status != service.Success {
		slog.Debug("ReadBackupDestinationsFailure", slog.Any("serviceOutput", readBackupDestinationsServiceOutput))
		return nil
	}

	readResponseDto, assertOk := readBackupDestinationsServiceOutput.Body.(dto.ReadBackupDestinationsResponse)
	if !assertOk {
		slog.Debug("AssertBackupDestinationsResponseDtoFailure")
		return nil
	}

	for _, iDestinationEntity := range readResponseDto.Destinations {
		selectLabelValuePair := componentForm.SelectLabelValuePair{
			Label: iDestinationEntity.ReadDestinationName().String() +
				" (#" + iDestinationEntity.ReadDestinationId().String() +
				" - " + iDestinationEntity.ReadDestinationType().String() + ")",
			Value: iDestinationEntity.ReadDestinationId().String(),
		}
		selectLabelValuePairs = append(selectLabelValuePairs, selectLabelValuePair)
	}

	return selectLabelValuePairs
}
