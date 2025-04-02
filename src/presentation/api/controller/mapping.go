package apiController

import (
	"net/http"

	"github.com/goinfinite/ez/src/domain/valueObject"
	"github.com/goinfinite/ez/src/infra/db"
	apiHelper "github.com/goinfinite/ez/src/presentation/api/helper"
	"github.com/goinfinite/ez/src/presentation/service"
	"github.com/labstack/echo/v4"
)

type MappingController struct {
	mappingService *service.MappingService
}

func NewMappingController(
	persistentDbSvc *db.PersistentDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) *MappingController {
	return &MappingController{
		mappingService: service.NewMappingService(persistentDbSvc, trailDbSvc),
	}
}

// ReadMappings	 godoc
// @Summary      ReadMappings
// @Description  List mappings.
// @Tags         mapping
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {array} entity.Mapping
// @Router       /v1/mapping/ [get]
func (controller *MappingController) Read(c echo.Context) error {
	return apiHelper.ServiceResponseWrapper(c, controller.mappingService.Read())
}

// CreateMapping	 godoc
// @Summary      CreateNewMapping
// @Description  Create a new mapping.
// @Tags         mapping
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createMappingDto 	  body    dto.CreateMapping  true  "NewMapping"
// @Success      201 {object} object{} "MappingCreated"
// @Router       /v1/mapping/ [post]
func (controller *MappingController) Create(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	if requestBody["containerId"] != nil {
		requestBody["containerIds"] = requestBody["containerId"]
	}

	containerIds := []valueObject.ContainerId{}
	if requestBody["containerIds"] != nil {
		_, isContainerIdsString := requestBody["containerIds"].(string)
		if isContainerIdsString {
			requestBody["containerIds"] = []interface{}{requestBody["containerIds"]}
		}

		containerIdsSlice, assertOk := requestBody["containerIds"].([]interface{})
		if !assertOk {
			return apiHelper.ResponseWrapper(
				c, http.StatusBadRequest, "ContainerIdsMustBeArray",
			)
		}

		for _, rawContainerId := range containerIdsSlice {
			containerId, err := valueObject.NewContainerId(rawContainerId)
			if err != nil {
				return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
			}
			containerIds = append(containerIds, containerId)
		}
	}
	requestBody["containerIds"] = containerIds

	return apiHelper.ServiceResponseWrapper(
		c, controller.mappingService.Create(requestBody),
	)
}

// DeleteMapping godoc
// @Summary      DeleteMapping
// @Description  Delete a mapping.
// @Tags         mapping
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        mappingId 	  path   string  true  "MappingId"
// @Success      200 {object} object{} "MappingDeleted"
// @Router       /v1/mapping/{mappingId}/ [delete]
// @Router       /v1/mapping/{accountId}/{mappingId}/ [delete]
func (controller *MappingController) Delete(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.mappingService.Delete(requestBody),
	)
}

// CreateMappingTarget godoc
// @Summary      CreateMappingTarget
// @Description  Create a new mapping target.
// @Tags         mapping
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createMappingTargetDto 	  body    dto.CreateMappingTarget  true  "NewMappingTarget"
// @Success      201 {object} object{} "MappingTargetCreated"
// @Router       /v1/mapping/target/ [post]
func (controller *MappingController) CreateTarget(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.mappingService.CreateTarget(requestBody),
	)
}

// DeleteMappingTarget godoc
// @Summary      DeleteMappingTarget
// @Description  Delete a mapping target.
// @Tags         mapping
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        mappingId 	  path   string  true  "MappingId"
// @Param        targetId 	  path   string  true  "TargetId"
// @Success      200 {object} object{} "MappingTargetDeleted"
// @Router       /v1/mapping/{mappingId}/target/{targetId}/ [delete]
// @Router       /v1/mapping/{accountId}/{mappingId}/target/{targetId}/ [delete]
func (controller *MappingController) DeleteTarget(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.mappingService.DeleteTarget(requestBody),
	)
}
