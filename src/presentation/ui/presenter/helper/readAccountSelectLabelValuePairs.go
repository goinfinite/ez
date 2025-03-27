package presenterHelper

import (
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/infra/db"
	"github.com/goinfinite/ez/src/presentation/service"
	componentForm "github.com/goinfinite/ez/src/presentation/ui/component/form"
)

func ReadAccountSelectLabelValuePairs(
	persistentDbSvc *db.PersistentDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) []componentForm.SelectLabelValuePair {
	selectLabelValuePairs := []componentForm.SelectLabelValuePair{}

	accountService := service.NewAccountService(persistentDbSvc, trailDbSvc)

	readAccountsServiceOutput := accountService.Read(map[string]any{})
	if readAccountsServiceOutput.Status != service.Success {
		slog.Debug("ReadAccountsFailure", slog.Any("serviceOutput", readAccountsServiceOutput))
		return nil
	}

	readAccountsRequestDto, assertOk := readAccountsServiceOutput.Body.(dto.ReadAccountsResponse)
	if !assertOk {
		slog.Debug("AssertAccountsFailure")
		return nil
	}

	for _, accountEntity := range readAccountsRequestDto.Accounts {
		accountIdStr := accountEntity.Id.String()
		selectLabelValuePair := componentForm.SelectLabelValuePair{
			Label: accountEntity.Username.String() + " (#" + accountIdStr + ")",
			Value: accountIdStr,
		}
		selectLabelValuePairs = append(selectLabelValuePairs, selectLabelValuePair)
	}

	return selectLabelValuePairs
}
