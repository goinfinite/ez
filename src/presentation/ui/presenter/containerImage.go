package presenter

import (
	"github.com/labstack/echo/v4"

	"net/http"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
	"github.com/speedianet/control/src/presentation/service"
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

	accountIdUsernameMap := map[valueObject.AccountId]valueObject.Username{}
	for _, accountEntity := range accountEntities {
		accountIdUsernameMap[accountEntity.Id] = accountEntity.Username
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

	pageContent := page.ContainerImageIndex(
		imageEntities, archiveFileEntities, accountIdUsernameMap,
		containerEntities, containerProfileEntities, accountEntities,
	)
	return uiHelper.Render(c, pageContent, http.StatusOK)
}
