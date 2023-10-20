package apiController

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/useCase"
	"github.com/speedianet/sfm/src/domain/valueObject"
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
// @Param        addContainerDto 	  body    dto.AddContainer  true  "NewContainer (Only accId, hostname and imgAddr are required.)"
// @Success      201 {object} object{} "ContainerCreated"
// @Router       /container/ [post]
func AddContainerController(c echo.Context) error {
	requiredParams := []string{"accId", "hostname", "imgAddr"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	accId := valueObject.NewAccountIdPanic(requestBody["accId"])
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

	var resourceProfileIdPtr *valueObject.ResourceProfileId
	if requestBody["resourceProfileId"] != nil {
		resourceProfileId := valueObject.NewResourceProfileIdPanic(
			requestBody["resourceProfileId"],
		)
		resourceProfileIdPtr = &resourceProfileId
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
		resourceProfileIdPtr,
		envs,
	)

	dbSvc := c.Get("dbSvc").(*gorm.DB)
	containerCmdRepo := infra.ContainerCmdRepo{}
	accQueryRepo := infra.NewAccQueryRepo(dbSvc)
	accCmdRepo := infra.NewAccCmdRepo(dbSvc)
	resourceProfileQueryRepo := infra.NewResourceProfileQueryRepo(dbSvc)

	err := useCase.AddContainer(
		containerCmdRepo,
		accQueryRepo,
		accCmdRepo,
		resourceProfileQueryRepo,
		addContainerDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusCreated, "ContainerCreated")
}

// UpdateContainer godoc
// @Summary      UpdateContainer
// @Description  Update an container.
// @Tags         container
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        updateContainerDto 	  body dto.UpdateContainer  true  "UpdateContainer (Only accId and containerId are required.)"
// @Success      200 {object} object{} "ContainerUpdated message or NewKeyString"
// @Router       /container/ [put]
func UpdateContainerController(c echo.Context) error {
	requiredParams := []string{"accId", "containerId"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	accId := valueObject.NewAccountIdPanic(requestBody["accId"])
	containerId := valueObject.NewContainerIdPanic(
		requestBody["containerId"].(string),
	)

	var containerStatusPtr *bool
	if requestBody["status"] != nil {
		containerStatus, err := strconv.ParseBool(requestBody["status"].(string))
		if err != nil {
			return apiHelper.ResponseWrapper(
				c, http.StatusBadRequest, "InvalidContainerStatus",
			)
		}
		containerStatusPtr = &containerStatus
	}

	var resourceProfileIdPtr *valueObject.ResourceProfileId
	if requestBody["resourceProfileId"] != nil {
		resourceProfileId := valueObject.NewResourceProfileIdPanic(
			requestBody["resourceProfileId"],
		)
		resourceProfileIdPtr = &resourceProfileId
	}

	updateContainerDto := dto.NewUpdateContainer(
		accId,
		containerId,
		containerStatusPtr,
		resourceProfileIdPtr,
	)

	dbSvc := c.Get("dbSvc").(*gorm.DB)
	containerQueryRepo := infra.NewContainerQueryRepo(dbSvc)
	containerCmdRepo := infra.ContainerCmdRepo{}
	accQueryRepo := infra.NewAccQueryRepo(dbSvc)
	accCmdRepo := infra.NewAccCmdRepo(dbSvc)
	resourceProfileQueryRepo := infra.NewResourceProfileQueryRepo(dbSvc)

	err := useCase.UpdateContainer(
		containerQueryRepo,
		containerCmdRepo,
		accQueryRepo,
		accCmdRepo,
		resourceProfileQueryRepo,
		updateContainerDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(
			c, http.StatusInternalServerError, err.Error(),
		)
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "ContainerUpdated")
}
