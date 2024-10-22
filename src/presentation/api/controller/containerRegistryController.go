package apiController

import (
	"net/http"

	"github.com/goinfinite/ez/src/domain/useCase"
	"github.com/goinfinite/ez/src/domain/valueObject"
	"github.com/goinfinite/ez/src/infra"
	"github.com/goinfinite/ez/src/infra/db"
	apiHelper "github.com/goinfinite/ez/src/presentation/api/helper"
	"github.com/labstack/echo/v4"
)

// GetContainerRegistryImages	 godoc
// @Summary      GetContainerRegistryImages
// @Description  Get container registry images.
// @Tags         containerRegistry
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        name    query     string  false  "ImageName"
// @Success      200 {array} entity.RegistryImage
// @Router       /v1/container/registry/image/ [get]
func GetContainerRegistryImagesController(c echo.Context) error {
	persistentDbSvc := c.Get("persistentDbSvc").(*db.PersistentDatabaseService)

	rawImageName := c.QueryParam("name")
	var imageNamePtr *valueObject.RegistryImageName
	if rawImageName != "" {
		imageName, err := valueObject.NewRegistryImageName(rawImageName)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
		imageNamePtr = &imageName
	}

	containerRegistryQueryRepo := infra.NewContainerRegistryQueryRepo(persistentDbSvc)
	imagesList, err := useCase.ReadRegistryImages(
		containerRegistryQueryRepo,
		imageNamePtr,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, imagesList)
}

// GetContainerRegistryTaggedImage	 godoc
// @Summary      GetContainerRegistryTaggedImage
// @Description  Get container registry tagged image.
// @Tags         containerRegistry
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        address    query     string  true  "ImageAddress"
// @Success      200 {object} entity.RegistryTaggedImage
// @Router       /v1/container/registry/image/tagged/ [get]
func GetContainerRegistryTaggedImageController(c echo.Context) error {
	persistentDbSvc := c.Get("persistentDbSvc").(*db.PersistentDatabaseService)

	imageAddressStr := c.QueryParam("address")
	imageAddress := valueObject.NewContainerImageAddressPanic(imageAddressStr)

	containerRegistryQueryRepo := infra.NewContainerRegistryQueryRepo(persistentDbSvc)
	taggedImage, err := useCase.ReadRegistryTaggedImage(
		containerRegistryQueryRepo,
		imageAddress,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, taggedImage)
}
