package presenterHelper

import (
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/infra/db"
	"github.com/goinfinite/ez/src/presentation/service"
	componentForm "github.com/goinfinite/ez/src/presentation/ui/component/form"
)

func ReadContainerSelectLabelValuePairs(
	persistentDbSvc *db.PersistentDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) []componentForm.SelectLabelValuePair {
	selectLabelValuePairs := []componentForm.SelectLabelValuePair{}

	containerService := service.NewContainerService(persistentDbSvc, trailDbSvc)

	readContainersServiceOutput := containerService.Read(map[string]interface{}{
		"itemsPerPage": 1000,
	})
	if readContainersServiceOutput.Status != service.Success {
		slog.Debug("ReadContainersFailure", slog.Any("serviceOutput", readContainersServiceOutput))
		return nil
	}

	readResponseDto, assertOk := readContainersServiceOutput.Body.(dto.ReadContainersResponse)
	if !assertOk {
		slog.Debug("AssertContainersResponseDtoFailure")
		return nil
	}

	for _, containerEntity := range readResponseDto.Containers {
		selectLabelValuePair := componentForm.SelectLabelValuePair{
			Label: containerEntity.Hostname.String() + " (#" + containerEntity.Id.String() + ")",
			Value: containerEntity.Id.String(),
		}
		selectLabelValuePairs = append(selectLabelValuePairs, selectLabelValuePair)
	}

	return selectLabelValuePairs
}
