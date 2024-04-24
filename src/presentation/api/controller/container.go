package apiController

import (
	"log"
	"strconv"

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
	return apiHelper.NewResponseWrapper(c, controller.containerService.Read())
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
	return apiHelper.NewResponseWrapper(
		c,
		controller.containerService.ReadWithMetrics(),
	)
}

func (controller *ContainerController) parsePortBindings(
	rawPortBindings []interface{},
) []valueObject.PortBinding {
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
func (controller *ContainerController) Create(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	portBindings := []valueObject.PortBinding{}
	if requestBody["portBindings"] != nil {
		_, isMapStringInterface := requestBody["portBindings"].(map[string]interface{})
		if isMapStringInterface {
			requestBody["portBindings"] = []interface{}{requestBody["portBindings"]}
		}

		portBindings = controller.parsePortBindings(
			requestBody["portBindings"].([]interface{}),
		)
	}

	envs := []valueObject.ContainerEnv{}
	if requestBody["envs"] != nil {
		envs = controller.parseContainerEnvs(requestBody["envs"].([]interface{}))
	}

	requestBody["portBindings"] = portBindings
	requestBody["envs"] = envs

	return apiHelper.NewResponseWrapper(
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

	return apiHelper.NewResponseWrapper(
		c,
		controller.containerService.Update(requestBody),
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

	return apiHelper.NewResponseWrapper(
		c,
		controller.containerService.Delete(requestBody),
	)
}
