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
// @Router       /v1/mapping/ [get]
func GetMappingsController(c echo.Context) error {
	persistentDbSvc := c.Get("persistentDbSvc").(*db.PersistentDatabaseService)
	mappingQueryRepo := infra.NewMappingQueryRepo(persistentDbSvc)
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
// @Router       /v1/mapping/ [post]
func AddMappingController(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	requiredParams := []string{"accountId", "publicPort", "containerIds"}
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

	containerIds := []valueObject.ContainerId{}
	for _, targetStr := range requestBody["containerIds"].([]interface{}) {
		containerId, err := valueObject.NewContainerId(targetStr.(string))
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
		containerIds = append(containerIds, containerId)
	}

	addMappingDto := dto.NewAddMapping(
		accId,
		hostnamePtr,
		publicPort,
		protocol,
		containerIds,
	)

	persistentDbSvc := c.Get("persistentDbSvc").(*db.PersistentDatabaseService)
	mappingQueryRepo := infra.NewMappingQueryRepo(persistentDbSvc)
	mappingCmdRepo := infra.NewMappingCmdRepo(persistentDbSvc)
	containerQueryRepo := infra.NewContainerQueryRepo(persistentDbSvc)

	err = useCase.AddMapping(
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
// @Router       /v1/mapping/{mappingId}/ [delete]
func DeleteMappingController(c echo.Context) error {
	mappingId := valueObject.NewMappingIdPanic(c.Param("mappingId"))

	persistentDbSvc := c.Get("persistentDbSvc").(*db.PersistentDatabaseService)
	mappingQueryRepo := infra.NewMappingQueryRepo(persistentDbSvc)
	mappingCmdRepo := infra.NewMappingCmdRepo(persistentDbSvc)

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
// @Router       /v1/mapping/target/ [post]
func AddMappingTargetController(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	requiredParams := []string{"mappingId", "containerId"}
	apiHelper.CheckMissingParams(requestBody, requiredParams)

	mappingId := valueObject.NewMappingIdPanic(requestBody["mappingId"])
	containerId := valueObject.NewContainerIdPanic(requestBody["containerId"].(string))

	addTargetDto := dto.NewAddMappingTarget(
		mappingId,
		containerId,
	)

	persistentDbSvc := c.Get("persistentDbSvc").(*db.PersistentDatabaseService)
	mappingQueryRepo := infra.NewMappingQueryRepo(persistentDbSvc)
	mappingCmdRepo := infra.NewMappingCmdRepo(persistentDbSvc)
	containerQueryRepo := infra.NewContainerQueryRepo(persistentDbSvc)

	err = useCase.AddMappingTarget(
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
// @Router       /v1/mapping/{mappingId}/target/{targetId}/ [delete]
func DeleteMappingTargetController(c echo.Context) error {
	targetId := valueObject.NewMappingTargetIdPanic(c.Param("targetId"))

	persistentDbSvc := c.Get("persistentDbSvc").(*db.PersistentDatabaseService)
	mappingQueryRepo := infra.NewMappingQueryRepo(persistentDbSvc)
	mappingCmdRepo := infra.NewMappingCmdRepo(persistentDbSvc)

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
