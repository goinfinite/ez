package presenterHelper

import (
	"log/slog"

	"github.com/goinfinite/ez/src/domain/entity"
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

	readAccountsServiceOutput := accountService.Read()
	if readAccountsServiceOutput.Status != service.Success {
		slog.Debug("ReadAccountsFailure", slog.Any("serviceOutput", readAccountsServiceOutput))
		return nil
	}

	accountEntities, assertOk := readAccountsServiceOutput.Body.([]entity.Account)
	if !assertOk {
		slog.Debug("AssertAccountsFailure")
		return nil
	}

	for _, accountEntity := range accountEntities {
		selectLabelValuePair := componentForm.SelectLabelValuePair{
			Label: accountEntity.Username.String() + " (#" + accountEntity.Id.String() + ")",
			Value: accountEntity.Id.String(),
		}
		selectLabelValuePairs = append(selectLabelValuePairs, selectLabelValuePair)
	}

	return selectLabelValuePairs
}
