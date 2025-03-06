package presenterHelper

import (
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
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
		var destinationName valueObject.BackupDestinationName
		var destinationId valueObject.BackupDestinationId
		var destinationType valueObject.BackupDestinationType

		switch destinationEntity := iDestinationEntity.(type) {
		case entity.BackupDestinationLocal:
			destinationName = destinationEntity.DestinationName
			destinationId = destinationEntity.DestinationId
			destinationType = destinationEntity.DestinationType
		case entity.BackupDestinationObjectStorage:
			destinationName = destinationEntity.DestinationName
			destinationId = destinationEntity.DestinationId
			destinationType = destinationEntity.DestinationType
		case entity.BackupDestinationRemoteHost:
			destinationName = destinationEntity.DestinationName
			destinationId = destinationEntity.DestinationId
			destinationType = destinationEntity.DestinationType
		}
		selectLabelValuePair := componentForm.SelectLabelValuePair{
			Label: destinationName.String() +
				" (#" + destinationId.String() +
				" - " + destinationType.String() + ")",
			Value: destinationId.String(),
		}
		selectLabelValuePairs = append(selectLabelValuePairs, selectLabelValuePair)
	}

	return selectLabelValuePairs
}
