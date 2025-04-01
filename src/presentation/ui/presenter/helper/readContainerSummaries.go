package presenterHelper

import (
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/infra/db"
	"github.com/goinfinite/ez/src/presentation/service"
	componentContainer "github.com/goinfinite/ez/src/presentation/ui/component/container"
)

func ReadContainerSummaries(
	persistentDbSvc *db.PersistentDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) []componentContainer.ContainerSummary {
	containerService := service.NewContainerService(persistentDbSvc, trailDbSvc)

	readContainersServiceOutput := containerService.Read(map[string]any{})
	if readContainersServiceOutput.Status != service.Success {
		slog.Debug("ReadContainersFailure")
		return nil
	}

	containersResponseDto, assertOk := readContainersServiceOutput.Body.(dto.ReadContainersResponse)
	if !assertOk {
		slog.Debug("AssertContainersResponseFailure")
		return nil
	}

	containerProfileService := service.NewContainerProfileService(persistentDbSvc, trailDbSvc)

	readContainerProfilesServiceOutput := containerProfileService.Read()
	if readContainerProfilesServiceOutput.Status != service.Success {
		slog.Debug("ReadContainerProfilesFailure")
		return nil
	}

	profileEntities, assertOk := readContainerProfilesServiceOutput.Body.([]entity.ContainerProfile)
	if !assertOk {
		slog.Debug("AssertContainerProfilesFailure")
		return nil
	}

	accountService := service.NewAccountService(persistentDbSvc, trailDbSvc)

	readAccountsServiceOutput := accountService.Read(map[string]any{})
	if readAccountsServiceOutput.Status != service.Success {
		slog.Debug("ReadAccountsFailure")
		return nil
	}

	readAccountsRequestDto, assertOk := readAccountsServiceOutput.Body.(dto.ReadAccountsResponse)
	if !assertOk {
		slog.Debug("AssertAccountsFailure")
		return nil
	}

	return componentContainer.NewContainerSummaries(
		containersResponseDto.Containers, profileEntities,
		readAccountsRequestDto.Accounts,
	)
}
