package presenter

import (
	"github.com/labstack/echo/v4"

	"net/http"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
	"github.com/speedianet/control/src/presentation/service"
	componentContainer "github.com/speedianet/control/src/presentation/ui/component/container"
	componentForm "github.com/speedianet/control/src/presentation/ui/component/form"
	uiHelper "github.com/speedianet/control/src/presentation/ui/helper"
	"github.com/speedianet/control/src/presentation/ui/page"
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

func transformContainerSummariesIntoSearchableSelectItems(
	containerEntities []entity.Container,
	profileEntities []entity.ContainerProfile,
	accountIdEntityMap map[valueObject.AccountId]entity.Account,
) []componentForm.SearchableSelectItem {
	searchableSelectItems := []componentForm.SearchableSelectItem{}

	containerIdEntityMap := map[valueObject.ContainerId]entity.Container{}
	for _, containerEntity := range containerEntities {
		containerIdEntityMap[containerEntity.Id] = containerEntity
	}

	profileIdEntityMap := map[valueObject.ContainerProfileId]entity.ContainerProfile{}
	for _, profileEntity := range profileEntities {
		profileIdEntityMap[profileEntity.Id] = profileEntity
	}

	containerSummaries := componentContainer.NewContainerSummariesWithMaps(
		containerIdEntityMap, profileIdEntityMap, accountIdEntityMap,
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
		return nil
	}

	imageEntities, assertOk := readImagesServiceOutput.Body.([]entity.ContainerImage)
	if !assertOk {
		return nil
	}

	readArchiveFilesServiceOutput := containerImageService.ReadArchiveFiles(&c.Request().Host)
	if readArchiveFilesServiceOutput.Status != service.Success {
		return nil
	}

	archiveFileEntities, assertOk := readArchiveFilesServiceOutput.Body.([]entity.ContainerImageArchiveFile)
	if !assertOk {
		return nil
	}

	accountService := service.NewAccountService(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)

	readAccountsServiceOutput := accountService.Read()
	if readAccountsServiceOutput.Status != service.Success {
		return nil
	}

	accountEntities, assertOk := readAccountsServiceOutput.Body.([]entity.Account)
	if !assertOk {
		return nil
	}

	accountIdEntityMap := map[valueObject.AccountId]entity.Account{}
	for _, accountEntity := range accountEntities {
		accountIdEntityMap[accountEntity.Id] = accountEntity
	}

	containerService := service.NewContainerService(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)

	readContainersServiceOutput := containerService.Read()
	if readContainersServiceOutput.Status != service.Success {
		return nil
	}

	containerEntities, assertOk := readContainersServiceOutput.Body.([]entity.Container)
	if !assertOk {
		return nil
	}

	containerProfileService := service.NewContainerProfileService(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)

	readContainerProfilesServiceOutput := containerProfileService.Read()
	if readContainerProfilesServiceOutput.Status != service.Success {
		return nil
	}

	containerProfileEntities, assertOk := readContainerProfilesServiceOutput.Body.([]entity.ContainerProfile)
	if !assertOk {
		return nil
	}

	accountsSelectPairs := transformAccountMapIntoSelectLabelValuePair(accountIdEntityMap)
	containerSummariesSearchableItems := transformContainerSummariesIntoSearchableSelectItems(
		containerEntities, containerProfileEntities, accountIdEntityMap,
	)

	pageContent := page.ContainerImageIndex(
		imageEntities, archiveFileEntities, accountIdEntityMap,
		accountsSelectPairs, containerSummariesSearchableItems,
	)
	return uiHelper.Render(c, pageContent, http.StatusOK)
}
