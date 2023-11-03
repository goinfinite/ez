package apiController

import (
	"net/http"
	"strconv"

	"github.com/goinfinite/fleet/src/domain/dto"
	"github.com/goinfinite/fleet/src/domain/useCase"
	"github.com/goinfinite/fleet/src/domain/valueObject"
	"github.com/goinfinite/fleet/src/infra"
	"github.com/goinfinite/fleet/src/infra/db"
	apiHelper "github.com/goinfinite/fleet/src/presentation/api/helper"
	"github.com/labstack/echo/v4"
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
	dbSvc := c.Get("dbSvc").(*db.DatabaseService)
	containerQueryRepo := infra.NewContainerQueryRepo(dbSvc)
	containerList, err := useCase.GetContainers(containerQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, containerList)
}

// GetContainersWithUsage	 godoc
// @Summary      GetContainersWithUsage
// @Description  List containers with usage.
// @Tags         container
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {array} dto.ContainerWithUsage
// @Router       /container/with-usage/ [get]
func GetContainersWithUsageController(c echo.Context) error {
	dbSvc := c.Get("dbSvc").(*db.DatabaseService)
	containerQueryRepo := infra.NewContainerQueryRepo(dbSvc)
	containerList, err := useCase.GetContainersWithUsage(containerQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, containerList)
}

func parsePortBindings(portBindings []interface{}) []valueObject.PortBinding {
	portBindingsList := []valueObject.PortBinding{}
	for _, portBinding := range portBindings {
		portBindingMap := portBinding.(map[string]interface{})
		protocol, err := valueObject.NewNetworkProtocol(
			portBindingMap["protocol"].(string),
		)
		if err != nil {
			continue
		}

		containerPort, err := valueObject.NewNetworkPort(
			portBindingMap["containerPort"].(string),
		)
		if err != nil {
			continue
		}

		hostPort, err := valueObject.NewNetworkPort(
			portBindingMap["hostPort"].(string),
		)
		if err != nil {
			continue
		}

		newPortBinding := valueObject.NewPortBinding(
			protocol,
			containerPort,
			hostPort,
		)
		portBindingsList = append(portBindingsList, newPortBinding)
	}

	return portBindingsList
}

func parseContainerEnvs(envs []interface{}) []valueObject.ContainerEnv {
	containerEnvs := []valueObject.ContainerEnv{}
	for _, env := range envs {
		newEnv, err := valueObject.NewContainerEnv(env.(string))
		if err != nil {
			continue
		}
		containerEnvs = append(containerEnvs, newEnv)
	}

	return containerEnvs
}

// AddContainer	 godoc
// @Summary      AddNewContainer
// @Description  Add a new container.
// @Tags         container
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        addContainerDto 	  body    dto.AddContainer  true  "NewContainer (Only accountId, hostname and imgAddr are required.)"
// @Success      201 {object} object{} "ContainerCreated"
// @Router       /container/ [post]
func AddContainerController(c echo.Context) error {
	requiredParams := []string{"accountId", "hostname", "imgAddr"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	accId := valueObject.NewAccountIdPanic(requestBody["accountId"])
	hostname := valueObject.NewFqdnPanic(requestBody["hostname"].(string))
	imgAddr := valueObject.NewContainerImgAddressPanic(
		requestBody["imgAddr"].(string),
	)

	portBindings := []valueObject.PortBinding{}
	if requestBody["portBindings"] != nil {
		portBindings = parsePortBindings(requestBody["portBindings"].([]interface{}))
	}

	var restartPolicyPtr *valueObject.ContainerRestartPolicy
	if requestBody["restartPolicy"] != nil {
		restartPolicy := valueObject.NewContainerRestartPolicyPanic(
			requestBody["restartPolicy"].(string),
		)
		restartPolicyPtr = &restartPolicy
	}

	var entrypointPtr *valueObject.ContainerEntrypoint
	if requestBody["entrypoint"] != nil {
		entrypoint := valueObject.NewContainerEntrypointPanic(
			requestBody["entrypoint"].(string),
		)
		entrypointPtr = &entrypoint
	}

	var profileIdPtr *valueObject.ContainerProfileId
	if requestBody["profileId"] != nil {
		profileId := valueObject.NewContainerProfileIdPanic(
			requestBody["profileId"],
		)
		profileIdPtr = &profileId
	}

	envs := []valueObject.ContainerEnv{}
	if requestBody["envs"] != nil {
		envs = parseContainerEnvs(requestBody["envs"].([]interface{}))
	}

	addContainerDto := dto.NewAddContainer(
		accId,
		hostname,
		imgAddr,
		portBindings,
		restartPolicyPtr,
		entrypointPtr,
		profileIdPtr,
		envs,
	)

	dbSvc := c.Get("dbSvc").(*db.DatabaseService)
	containerCmdRepo := infra.NewContainerCmdRepo(dbSvc)
	accQueryRepo := infra.NewAccQueryRepo(dbSvc)
	accCmdRepo := infra.NewAccCmdRepo(dbSvc)
	containerProfileQueryRepo := infra.NewContainerProfileQueryRepo(dbSvc)

	err := useCase.AddContainer(
		containerCmdRepo,
		accQueryRepo,
		accCmdRepo,
		containerProfileQueryRepo,
		addContainerDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusCreated, "ContainerCreated")
}

// UpdateContainer godoc
// @Summary      UpdateContainer
// @Description  Update a container.
// @Tags         container
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        updateContainerDto 	  body dto.UpdateContainer  true  "UpdateContainer (Only accountId and containerId are required.)"
// @Success      200 {object} object{} "ContainerUpdated message or NewKeyString"
// @Router       /container/ [put]
func UpdateContainerController(c echo.Context) error {
	requiredParams := []string{"accountId", "containerId"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	accId := valueObject.NewAccountIdPanic(requestBody["accountId"])
	containerId := valueObject.NewContainerIdPanic(
		requestBody["containerId"].(string),
	)

	var containerStatusPtr *bool
	if requestBody["status"] != nil {
		containerStatus, assertOk := requestBody["status"].(bool)
		if !assertOk {
			var err error
			containerStatus, err = strconv.ParseBool(requestBody["status"].(string))
			if err != nil {
				return apiHelper.ResponseWrapper(
					c, http.StatusBadRequest, "InvalidContainerStatus",
				)
			}
		}
		containerStatusPtr = &containerStatus
	}

	var profileIdPtr *valueObject.ContainerProfileId
	if requestBody["profileId"] != nil {
		profileId := valueObject.NewContainerProfileIdPanic(
			requestBody["profileId"],
		)
		profileIdPtr = &profileId
	}

	updateContainerDto := dto.NewUpdateContainer(
		accId,
		containerId,
		containerStatusPtr,
		profileIdPtr,
	)

	dbSvc := c.Get("dbSvc").(*db.DatabaseService)
	containerQueryRepo := infra.NewContainerQueryRepo(dbSvc)
	containerCmdRepo := infra.NewContainerCmdRepo(dbSvc)
	accQueryRepo := infra.NewAccQueryRepo(dbSvc)
	accCmdRepo := infra.NewAccCmdRepo(dbSvc)
	containerProfileQueryRepo := infra.NewContainerProfileQueryRepo(dbSvc)

	err := useCase.UpdateContainer(
		containerQueryRepo,
		containerCmdRepo,
		accQueryRepo,
		accCmdRepo,
		containerProfileQueryRepo,
		updateContainerDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(
			c, http.StatusInternalServerError, err.Error(),
		)
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "ContainerUpdated")
}

// DeleteContainer godoc
// @Summary      DeleteContainer
// @Description  Delete a container.
// @Tags         container
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        accId 	  path   string  true  "AccountId"
// @Param        containerId 	  path   string  true  "ContainerId"
// @Success      200 {object} object{} "ContainerDeleted"
// @Router       /container/{accountId}/{containerId}/ [delete]
func DeleteContainerController(c echo.Context) error {
	accId := valueObject.NewAccountIdPanic(c.Param("accountId"))
	containerId := valueObject.NewContainerIdPanic(c.Param("containerId"))

	dbSvc := c.Get("dbSvc").(*db.DatabaseService)
	containerQueryRepo := infra.NewContainerQueryRepo(dbSvc)
	containerCmdRepo := infra.NewContainerCmdRepo(dbSvc)
	accCmdRepo := infra.NewAccCmdRepo(dbSvc)

	err := useCase.DeleteContainer(
		containerQueryRepo,
		containerCmdRepo,
		accCmdRepo,
		accId,
		containerId,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "ContainerDeleted")
}
