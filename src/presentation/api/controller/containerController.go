package apiController

import (
	"log"
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
// @Router       /v1/container/ [get]
func GetContainersController(c echo.Context) error {
	persistentDbSvc := c.Get("persistentDbSvc").(*db.PersistentDatabaseService)
	containerQueryRepo := infra.NewContainerQueryRepo(persistentDbSvc)
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
// @Router       /v1/container/metrics/ [get]
func GetContainersWithMetricsController(c echo.Context) error {
	persistentDbSvc := c.Get("persistentDbSvc").(*db.PersistentDatabaseService)
	containerQueryRepo := infra.NewContainerQueryRepo(persistentDbSvc)
	containerList, err := useCase.GetContainersWithMetrics(containerQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, containerList)
}

func parsePortBindings(rawPortBindings []interface{}) []valueObject.PortBinding {
	portBindings := []valueObject.PortBinding{}
	for rawPortBindingIndex, rawPortBinding := range rawPortBindings {
		errMsgSuffix := ": (item " + strconv.Itoa(rawPortBindingIndex) + ")"

		rawPortBindingMap, assertOk := rawPortBinding.(map[string]interface{})
		if !assertOk {
			log.Print("InvalidPortBindingStructure" + errMsgSuffix)
			continue
		}

		portBindingStr := ""

		rawServiceName, assertOk := rawPortBindingMap["serviceName"].(string)
		if assertOk {
			portBindingStr += rawServiceName
		}

		rawPublicPort, exists := rawPortBindingMap["publicPort"]
		if exists {
			publicPort, err := valueObject.NewNetworkPort(rawPublicPort)
			if err != nil {
				log.Print(err.Error() + errMsgSuffix)
				continue
			}
			portBindingStr += ":" + publicPort.String()
		}

		rawContainerPort, rawContainerPortExists := rawPortBindingMap["containerPort"]
		if rawContainerPortExists {
			containerPort, err := valueObject.NewNetworkPort(rawContainerPort)
			if err != nil {
				log.Print(err.Error() + errMsgSuffix)
				continue
			}
			portBindingStr += ":" + containerPort.String()
		}

		rawProtocol, assertOk := rawPortBindingMap["protocol"].(string)
		if assertOk && rawContainerPortExists {
			protocol, err := valueObject.NewNetworkProtocol(rawProtocol)
			if err != nil {
				log.Print(err.Error() + errMsgSuffix)
				continue
			}
			portBindingStr += "/" + protocol.String()
		}

		rawPrivatePort, exists := rawPortBindingMap["privatePort"]
		if exists {
			privatePort, err := valueObject.NewNetworkPort(rawPrivatePort)
			if err != nil {
				log.Print(err.Error() + errMsgSuffix)
				continue
			}
			portBindingStr += ":" + privatePort.String()
		}

		portBinding, err := valueObject.NewPortBindingFromString(portBindingStr)
		if err != nil {
			log.Print(err.Error() + errMsgSuffix)
			continue
		}

		portBindings = append(portBindings, portBinding...)
	}

	return portBindings
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
// @Router       /v1/container/ [post]
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

	persistentDbSvc := c.Get("persistentDbSvc").(*db.PersistentDatabaseService)
	containerQueryRepo := infra.NewContainerQueryRepo(persistentDbSvc)
	containerCmdRepo := infra.NewContainerCmdRepo(persistentDbSvc)
	accQueryRepo := infra.NewAccQueryRepo(persistentDbSvc)
	accCmdRepo := infra.NewAccCmdRepo(persistentDbSvc)
	containerProfileQueryRepo := infra.NewContainerProfileQueryRepo(persistentDbSvc)
	mappingQueryRepo := infra.NewMappingQueryRepo(persistentDbSvc)
	mappingCmdRepo := infra.NewMappingCmdRepo(persistentDbSvc)

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
// @Router       /v1/container/ [put]
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

	persistentDbSvc := c.Get("persistentDbSvc").(*db.PersistentDatabaseService)
	containerQueryRepo := infra.NewContainerQueryRepo(persistentDbSvc)
	containerCmdRepo := infra.NewContainerCmdRepo(persistentDbSvc)
	accQueryRepo := infra.NewAccQueryRepo(persistentDbSvc)
	accCmdRepo := infra.NewAccCmdRepo(persistentDbSvc)
	containerProfileQueryRepo := infra.NewContainerProfileQueryRepo(persistentDbSvc)

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
// @Router       /v1/container/{accountId}/{containerId}/ [delete]
func DeleteContainerController(c echo.Context) error {
	accId := valueObject.NewAccountIdPanic(c.Param("accountId"))
	containerId := valueObject.NewContainerIdPanic(c.Param("containerId"))

	persistentDbSvc := c.Get("persistentDbSvc").(*db.PersistentDatabaseService)
	containerQueryRepo := infra.NewContainerQueryRepo(persistentDbSvc)
	containerCmdRepo := infra.NewContainerCmdRepo(persistentDbSvc)
	accCmdRepo := infra.NewAccCmdRepo(persistentDbSvc)
	mappingQueryRepo := infra.NewMappingQueryRepo(persistentDbSvc)
	mappingCmdRepo := infra.NewMappingCmdRepo(persistentDbSvc)

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
