package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/sfm/src/domain/useCase"
	"github.com/speedianet/sfm/src/infra"
	apiHelper "github.com/speedianet/sfm/src/presentation/api/helper"
	"gorm.io/gorm"
)

// GetContainers	 godoc
// @Summary      GetContainers
// @Description  List containers.
// @Tags         container
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {array} entity.Container
// @Router       /container/ [get]
func GetContainersController(c echo.Context) error {
	containerQueryRepo := infra.NewContainerQueryRepo(c.Get("dbSvc").(*gorm.DB))
	containerList, err := useCase.GetContainers(containerQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, containerList)
}
