package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
	apiHelper "github.com/speedianet/control/src/presentation/api/helper"
)

// GetMappings	 godoc
// @Summary      GetMappings
// @Description  List mapping.
// @Tags         mapping
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {array} entity.Mapping
// @Router       /mapping/ [get]
func GetMappingsController(c echo.Context) error {
	dbSvc := c.Get("dbSvc").(*db.DatabaseService)
	mappingQueryRepo := infra.NewMappingQueryRepo(dbSvc)
	mappingList, err := useCase.GetMappings(mappingQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, mappingList)
}

// AddMapping	 godoc
// @Summary      AddNewMapping
// @Description  Add a new mapping.
// @Tags         mapping
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        addMappingDto 	  body    dto.AddMapping  true  "NewMapping"
// @Success      201 {object} object{} "MappingCreated"
// @Router       /mapping/ [post]
func AddMappingController(c echo.Context) error {
	requiredParams := []string{"accountId", "publicPort", "targets"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	accId := valueObject.NewAccountIdPanic(requestBody["accountId"])

	var hostnamePtr *valueObject.Fqdn
	if requestBody["hostname"] != nil {
		hostname := valueObject.NewFqdnPanic(requestBody["hostname"].(string))
		hostnamePtr = &hostname
	}

	publicPort := valueObject.NewNetworkPortPanic(requestBody["publicPort"])

	protocol := valueObject.GuessNetworkProtocolByPort(publicPort)
	if requestBody["protocol"] != nil {
		protocol = valueObject.NewNetworkProtocolPanic(requestBody["protocol"].(string))
	}

	targets := []valueObject.ContainerId{}
	for _, targetStr := range requestBody["targets"].([]interface{}) {
		containerId, err := valueObject.NewContainerId(targetStr.(string))
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
		targets = append(targets, containerId)
	}

	addMappingDto := dto.NewAddMapping(
		accId,
		hostnamePtr,
		publicPort,
		protocol,
		targets,
	)

	dbSvc := c.Get("dbSvc").(*db.DatabaseService)
	mappingQueryRepo := infra.NewMappingQueryRepo(dbSvc)
	mappingCmdRepo := infra.NewMappingCmdRepo(dbSvc)
	containerQueryRepo := infra.NewContainerQueryRepo(dbSvc)

	err := useCase.AddMapping(
		mappingQueryRepo,
		mappingCmdRepo,
		containerQueryRepo,
		addMappingDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusCreated, "MappingCreated")
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
// @Router       /mapping/{mappingId}/ [delete]
func DeleteMappingController(c echo.Context) error {
	mappingId := valueObject.NewMappingIdPanic(c.Param("mappingId"))

	dbSvc := c.Get("dbSvc").(*db.DatabaseService)
	mappingQueryRepo := infra.NewMappingQueryRepo(dbSvc)
	mappingCmdRepo := infra.NewMappingCmdRepo(dbSvc)

	err := useCase.DeleteMapping(
		mappingQueryRepo,
		mappingCmdRepo,
		mappingId,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "MappingDeleted")
}

// AddMappingTarget godoc
// @Summary      AddMappingTarget
// @Description  Add a new mapping target.
// @Tags         mapping
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        addMappingTargetDto 	  body    dto.AddMappingTarget  true  "NewMappingTarget"
// @Success      201 {object} object{} "MappingTargetAdded"
// @Router       /mapping/target/ [post]
func AddMappingTargetController(c echo.Context) error {
	requiredParams := []string{"mappingId", "target"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	mappingId := valueObject.NewMappingIdPanic(requestBody["mappingId"])
	target := valueObject.NewContainerIdPanic(requestBody["target"].(string))

	addTargetDto := dto.NewAddMappingTarget(
		mappingId,
		target,
	)

	dbSvc := c.Get("dbSvc").(*db.DatabaseService)
	mappingQueryRepo := infra.NewMappingQueryRepo(dbSvc)
	mappingCmdRepo := infra.NewMappingCmdRepo(dbSvc)
	containerQueryRepo := infra.NewContainerQueryRepo(dbSvc)

	err := useCase.AddMappingTarget(
		mappingQueryRepo,
		mappingCmdRepo,
		containerQueryRepo,
		addTargetDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusCreated, "MappingTargetAdded")
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
// @Router       /mapping/{mappingId}/target/{targetId}/ [delete]
func DeleteMappingTargetController(c echo.Context) error {
	targetId := valueObject.NewMappingTargetIdPanic(c.Param("targetId"))

	dbSvc := c.Get("dbSvc").(*db.DatabaseService)
	mappingQueryRepo := infra.NewMappingQueryRepo(dbSvc)
	mappingCmdRepo := infra.NewMappingCmdRepo(dbSvc)

	err := useCase.DeleteMappingTarget(
		mappingQueryRepo,
		mappingCmdRepo,
		targetId,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "MappingTargetDeleted")
}
