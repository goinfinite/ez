package apiController

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
	apiHelper "github.com/speedianet/control/src/presentation/api/helper"
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

// GetContainersWithMetrics	 godoc
// @Summary      GetContainersWithMetrics
// @Description  List containers with metrics.
// @Tags         container
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {array} dto.ContainerWithMetrics
// @Router       /container/metrics/ [get]
func GetContainersWithMetricsController(c echo.Context) error {
	dbSvc := c.Get("dbSvc").(*db.DatabaseService)
	containerQueryRepo := infra.NewContainerQueryRepo(dbSvc)
	containerList, err := useCase.GetContainersWithMetrics(containerQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, containerList)
}

func parsePortBindings(portBindings []interface{}) []valueObject.PortBinding {
	portBindingsList := []valueObject.PortBinding{}
	for _, portBinding := range portBindings {
		portBindingMap := portBinding.(map[string]interface{})

		publicPort, err := valueObject.NewNetworkPort(
			portBindingMap["publicPort"],
		)
		if err != nil {
			continue
		}

		containerPort := publicPort
		if portBindingMap["containerPort"] != nil {
			containerPort, err = valueObject.NewNetworkPort(
				portBindingMap["containerPort"],
			)
			if err != nil {
				continue
			}
		}

		protocol := valueObject.GuessNetworkProtocolByPort(publicPort)
		if portBindingMap["protocol"] != nil {
			protocol, err = valueObject.NewNetworkProtocol(
				portBindingMap["protocol"].(string),
			)
			if err != nil {
				continue
			}
		}

		var privatePortPtr *valueObject.NetworkPort
		if portBindingMap["privatePort"] != nil {
			privatePort, err := valueObject.NewNetworkPort(
				portBindingMap["privatePort"],
			)
			if err != nil {
				continue
			}
			privatePortPtr = &privatePort
		}

		newPortBinding := valueObject.NewPortBinding(
			publicPort,
			containerPort,
			protocol,
			privatePortPtr,
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
// @Param        addContainerDto 	  body    dto.AddContainer  true  "NewContainer (Only accountId, hostname and imageAddress are required.)<br />When specifying portBindings, only publicPort is required."
// @Success      201 {object} object{} "ContainerCreated"
// @Router       /container/ [post]
func AddContainerController(c echo.Context) error {
	requiredParams := []string{"accountId", "hostname"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	accId := valueObject.NewAccountIdPanic(requestBody["accountId"])
	hostname := valueObject.NewFqdnPanic(requestBody["hostname"].(string))

	imgAddrStr, assertOk := requestBody["imageAddress"].(string)
	if !assertOk {
		imgAddrStr, assertOk = requestBody["imgAddr"].(string)
		if !assertOk {
			return apiHelper.ResponseWrapper(
				c, http.StatusBadRequest, "MissingImageAddress",
			)
		}
	}
	imgAddr := valueObject.NewContainerImageAddressPanic(imgAddrStr)

	portBindings := []valueObject.PortBinding{}
	if requestBody["portBindings"] != nil {
		_, isMapStringInterface := requestBody["portBindings"].(map[string]interface{})
		if isMapStringInterface {
			requestBody["portBindings"] = []interface{}{requestBody["portBindings"]}
		}

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

	autoCreateMappings := true
	if requestBody["autoCreateMappings"] != nil {
		var assertOk bool
		autoCreateMappings, assertOk = requestBody["autoCreateMappings"].(bool)
		if !assertOk {
			var err error
			autoCreateMappings, err = strconv.ParseBool(
				requestBody["autoCreateMappings"].(string),
			)
			if err != nil {
				return apiHelper.ResponseWrapper(
					c, http.StatusBadRequest, "InvalidAutoCreateMappings",
				)
			}
		}
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
		autoCreateMappings,
	)

	dbSvc := c.Get("dbSvc").(*db.DatabaseService)
	containerQueryRepo := infra.NewContainerQueryRepo(dbSvc)
	containerCmdRepo := infra.NewContainerCmdRepo(dbSvc)
	accQueryRepo := infra.NewAccQueryRepo(dbSvc)
	accCmdRepo := infra.NewAccCmdRepo(dbSvc)
	containerProfileQueryRepo := infra.NewContainerProfileQueryRepo(dbSvc)
	mappingQueryRepo := infra.NewMappingQueryRepo(dbSvc)
	mappingCmdRepo := infra.NewMappingCmdRepo(dbSvc)

	err := useCase.AddContainer(
		containerQueryRepo,
		containerCmdRepo,
		accQueryRepo,
		accCmdRepo,
		containerProfileQueryRepo,
		mappingQueryRepo,
		mappingCmdRepo,
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
// @Param        accountId 	  path   string  true  "AccountId"
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
	mappingQueryRepo := infra.NewMappingQueryRepo(dbSvc)
	mappingCmdRepo := infra.NewMappingCmdRepo(dbSvc)

	err := useCase.DeleteContainer(
		containerQueryRepo,
		containerCmdRepo,
		accCmdRepo,
		mappingQueryRepo,
		mappingCmdRepo,
		accId,
		containerId,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "ContainerDeleted")
}
