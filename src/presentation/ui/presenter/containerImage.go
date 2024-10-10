package presenter

import (
	"log/slog"

	"github.com/labstack/echo/v4"

	"net/http"

	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
	"github.com/goinfinite/ez/src/infra/db"
	"github.com/goinfinite/ez/src/presentation/service"
	componentContainer "github.com/goinfinite/ez/src/presentation/ui/component/container"
	componentForm "github.com/goinfinite/ez/src/presentation/ui/component/form"
	uiHelper "github.com/goinfinite/ez/src/presentation/ui/helper"
	"github.com/goinfinite/ez/src/presentation/ui/page"
	presenterHelper "github.com/goinfinite/ez/src/presentation/ui/presenter/helper"
)

type ContainerImagePresenter struct {
	persistentDbSvc *db.PersistentDatabaseService
	trailDbSvc      *db.TrailDatabaseService
}

func NewContainerImagePresenter(
	persistentDbSvc *db.PersistentDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) *ContainerImagePresenter {
	return &ContainerImagePresenter{
		persistentDbSvc: persistentDbSvc,
		trailDbSvc:      trailDbSvc,
	}
}

func transformAccountMapIntoSelectLabelValuePair(
	accountIdEntityMap map[valueObject.AccountId]entity.Account,
) []componentForm.SelectLabelValuePair {
	selectLabelValuePairs := []componentForm.SelectLabelValuePair{}
	for accountId, accountEntity := range accountIdEntityMap {
		selectLabelValuePair := componentForm.SelectLabelValuePair{
			Label: accountEntity.Username.String(),
			Value: accountId.String(),
		}
		selectLabelValuePairs = append(selectLabelValuePairs, selectLabelValuePair)
	}
	return selectLabelValuePairs
}

func (presenter *ContainerImagePresenter) transformContainerSummariesIntoSearchableItems() []componentForm.SearchableSelectItem {
	searchableSelectItems := []componentForm.SearchableSelectItem{}

	containerSummaries := presenterHelper.ReadContainerSummaries(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)

	for _, containerSummary := range containerSummaries {
		searchableTextSerialized := containerSummary.JsonSerialize()
		htmlLabel := componentContainer.ContainerTaggedSummary(containerSummary)

		searchableSelectItem := componentForm.SearchableSelectItem{
			Label:          containerSummary.Hostname.String(),
			Value:          containerSummary.ContainerId.String(),
			SearchableText: &searchableTextSerialized,
			HtmlLabel:      &htmlLabel,
		}
		searchableSelectItems = append(searchableSelectItems, searchableSelectItem)
	}

	return searchableSelectItems
}

func (presenter *ContainerImagePresenter) Handler(c echo.Context) error {
	containerImageService := service.NewContainerImageService(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)

	readImagesServiceOutput := containerImageService.Read()
	if readImagesServiceOutput.Status != service.Success {
		slog.Debug("ReadImagesFailure")
		return nil
	}

	imageEntities, assertOk := readImagesServiceOutput.Body.([]entity.ContainerImage)
	if !assertOk {
		slog.Debug("AssertImagesFailure")
		return nil
	}

	readArchiveFilesServiceOutput := containerImageService.ReadArchiveFiles(&c.Request().Host)
	if readArchiveFilesServiceOutput.Status != service.Success {
		slog.Debug("ReadArchiveFilesFailure")
		return nil
	}

	archiveFileEntities, assertOk := readArchiveFilesServiceOutput.Body.([]entity.ContainerImageArchiveFile)
	if !assertOk {
		slog.Debug("AssertArchiveFilesFailure")
		return nil
	}

	accountService := service.NewAccountService(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)

	readAccountsServiceOutput := accountService.Read()
	if readAccountsServiceOutput.Status != service.Success {
		slog.Debug("ReadAccountsFailure")
		return nil
	}

	accountEntities, assertOk := readAccountsServiceOutput.Body.([]entity.Account)
	if !assertOk {
		slog.Debug("AssertAccountsFailure")
		return nil
	}

	accountIdEntityMap := map[valueObject.AccountId]entity.Account{}
	for _, accountEntity := range accountEntities {
		accountIdEntityMap[accountEntity.Id] = accountEntity
	}

	accountsSelectPairs := transformAccountMapIntoSelectLabelValuePair(accountIdEntityMap)
	containerSummariesSearchableItems := presenter.transformContainerSummariesIntoSearchableItems()

	pageContent := page.ContainerImageIndex(
		imageEntities, archiveFileEntities, accountIdEntityMap,
		accountsSelectPairs, containerSummariesSearchableItems,
	)
	return uiHelper.Render(c, pageContent, http.StatusOK)
}
