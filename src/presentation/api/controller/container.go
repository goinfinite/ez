package apiController

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
	apiHelper "github.com/speedianet/control/src/presentation/api/helper"
	"github.com/speedianet/control/src/presentation/service"
)

type ContainerController struct {
	containerService *service.ContainerService
}

func NewContainerController(
	persistentDbSvc *db.PersistentDatabaseService,
) *ContainerController {
	return &ContainerController{
		containerService: service.NewContainerService(persistentDbSvc),
	}
}

// ReadContainers	 godoc
// @Summary      ReadContainers
// @Description  List containers.
// @Tags         container
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {array} entity.Container
// @Router       /v1/container/ [get]
func (controller *ContainerController) Read(c echo.Context) error {
	return apiHelper.ServiceResponseWrapper(c, controller.containerService.Read())
}

// ReadContainersWithMetrics	 godoc
// @Summary      ReadContainersWithMetrics
// @Description  List containers with metrics.
// @Tags         container
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {array} dto.ContainerWithMetrics
// @Router       /v1/container/metrics/ [get]
func (controller *ContainerController) ReadWithMetrics(c echo.Context) error {
	return apiHelper.ServiceResponseWrapper(
		c, controller.containerService.ReadWithMetrics(),
	)
}

func (controller *ContainerController) parsePortBindings(
	rawPortBindings []interface{},
) []valueObject.PortBinding {
	portBindings := []valueObject.PortBinding{}
	for bindingIndex, rawPortBinding := range rawPortBindings {
		rawPortBindingMap, assertOk := rawPortBinding.(map[string]interface{})
		if !assertOk {
			log.Printf("[%d] InvalidPortBindingStructure", bindingIndex)
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
				log.Printf("[%d] %s", bindingIndex, err.Error())
				continue
			}
			portBindingStr += ":" + publicPort.String()
		}

		rawContainerPort, rawContainerPortExists := rawPortBindingMap["containerPort"]
		if rawContainerPortExists {
			containerPort, err := valueObject.NewNetworkPort(rawContainerPort)
			if err != nil {
				log.Printf("[%d] %s", bindingIndex, err.Error())
				continue
			}
			portBindingStr += ":" + containerPort.String()
		}

		rawProtocol, assertOk := rawPortBindingMap["protocol"].(string)
		if assertOk && rawContainerPortExists {
			protocol, err := valueObject.NewNetworkProtocol(rawProtocol)
			if err != nil {
				log.Printf("[%d] %s", bindingIndex, err.Error())
				continue
			}
			portBindingStr += "/" + protocol.String()
		}

		rawPrivatePort, exists := rawPortBindingMap["privatePort"]
		if exists {
			privatePort, err := valueObject.NewNetworkPort(rawPrivatePort)
			if err != nil {
				log.Printf("[%d] %s", bindingIndex, err.Error())
				continue
			}
			portBindingStr += ":" + privatePort.String()
		}

		portBinding, err := valueObject.NewPortBindingFromString(portBindingStr)
		if err != nil {
			log.Printf("[%d] %s", bindingIndex, err.Error())
			continue
		}

		portBindings = append(portBindings, portBinding...)
	}

	return portBindings
}

func (controller *ContainerController) parseContainerEnvs(
	envs []interface{},
) []valueObject.ContainerEnv {
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

// CreateContainer	 godoc
// @Summary      CreateNewContainer
// @Description  Create a new container.
// @Tags         container
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createContainerDto 	  body    dto.CreateContainer  true  "NewContainer (Only accountId, hostname and imageAddress are required.)<br />When specifying portBindings, only publicPort is required.<br />LaunchScript must be base64 encoded."
// @Success      201 {object} object{} "ContainerCreated"
// @Router       /v1/container/ [post]
func (controller *ContainerController) Create(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	if _, exists := requestBody["imageAddress"]; !exists {
		if _, exists = requestBody["imgAddr"]; exists {
			requestBody["imageAddress"] = requestBody["imgAddr"]
		}
	}

	if requestBody["portBindings"] != nil {
		_, isMapStringInterface := requestBody["portBindings"].(map[string]interface{})
		if isMapStringInterface {
			requestBody["portBindings"] = []interface{}{requestBody["portBindings"]}
		}

		portBindingsSlice, assertOk := requestBody["portBindings"].([]interface{})
		if !assertOk {
			return apiHelper.ResponseWrapper(
				c, http.StatusBadRequest, "PortBindingsMustBeArray",
			)
		}

		portBindings := controller.parsePortBindings(portBindingsSlice)
		requestBody["portBindings"] = portBindings
	}

	if requestBody["envs"] != nil {
		envs := controller.parseContainerEnvs(requestBody["envs"].([]interface{}))
		requestBody["envs"] = envs
	}

	if requestBody["launchScript"] != nil {
		scriptEncodedContent, err := valueObject.NewEncodedContent(
			requestBody["launchScript"].(string),
		)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}

		scriptDecodedContent, err := scriptEncodedContent.GetDecoded()
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}

		launchScript, err := valueObject.NewLaunchScript(scriptDecodedContent)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}

		requestBody["launchScript"] = launchScript
	}

	return apiHelper.ServiceResponseWrapper(
		c,
		controller.containerService.Create(requestBody),
	)
}

// UpdateContainer godoc
// @Summary      UpdateContainer
// @Description  Update a container.
// @Tags         container
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        updateContainerDto 	  body dto.UpdateContainer  true  "UpdateContainer (Only accountId and containerId are required.)"
// @Success      200 {object} object{} "ContainerUpdated"
// @Router       /v1/container/ [put]
func (controller *ContainerController) Update(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.containerService.Update(requestBody),
	)
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
func (controller *ContainerController) Delete(c echo.Context) error {
	requestBody := map[string]interface{}{
		"accountId":   c.Param("accountId"),
		"containerId": c.Param("containerId"),
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.containerService.Delete(requestBody),
	)
}
