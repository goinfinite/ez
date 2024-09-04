package apiController

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/domain/valueObject"
	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
	"github.com/speedianet/control/src/infra/db"
	apiHelper "github.com/speedianet/control/src/presentation/api/helper"
	"github.com/speedianet/control/src/presentation/service"
)

type ContainerController struct {
	containerService *service.ContainerService
}

func NewContainerController(
	persistentDbSvc *db.PersistentDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) *ContainerController {
	return &ContainerController{
		containerService: service.NewContainerService(
			persistentDbSvc, trailDbSvc,
		),
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

// CreateContainerSessionToken	 godoc
// @Summary      CreateContainerSessionToken
// @Description  Creates a session token for the specified container and redirects to Speedia OS dashboard (if shouldRedirect is not false).
// @Tags         container
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        accountId	path	string	true	"AccountId"
// @Param        containerId	path	string	true	"ContainerId"
// @Param        shouldRedirect	query	bool	false	"ShouldRedirect (default/empty is true)"
// @Success      200 {object} valueObject.AccessTokenValue "If shouldRedirect is set to false, the updated session token is returned."
// @Success      302 {string} string "A redirect to Speedia OS dashboard (:1618/{containerId}/)."
// @Failure      500 {string} string "Container is not found, not running or isn't Speedia OS."
// @Router       /v1/container/session/{accountId}/{containerId}/ [get]
func (controller *ContainerController) CreateContainerSessionToken(c echo.Context) error {
	requestBody := map[string]interface{}{
		"accountId":         c.Param("accountId"),
		"containerId":       c.Param("containerId"),
		"sessionIpAddress":  c.RealIP(),
		"operatorAccountId": c.Get("accountId"),
		"operatorIpAddress": c.RealIP(),
	}

	serviceOutput := controller.containerService.CreateContainerSessionToken(requestBody)
	if serviceOutput.Status != service.Success {
		return apiHelper.ResponseWrapper(
			c, http.StatusInternalServerError, serviceOutput.Body,
		)
	}

	var err error
	shouldRedirect := true
	if rawShouldRedirect := c.QueryParam("shouldRedirect"); rawShouldRedirect != "" {
		shouldRedirect, err = voHelper.InterfaceToBool(rawShouldRedirect)
		if err != nil {
			shouldRedirect = false
		}
	}

	accessToken, assertOk := serviceOutput.Body.(valueObject.AccessTokenValue)
	if !assertOk {
		return apiHelper.ResponseWrapper(
			c, http.StatusInternalServerError, "InvalidAccessTokenValue",
		)
	}

	if !shouldRedirect {
		return apiHelper.ResponseWrapper(c, http.StatusOK, accessToken)
	}

	currentHost := c.Request().Host
	currentHost, _, err = net.SplitHostPort(currentHost)
	if err != nil {
		currentHost = c.Request().Host
	}

	redirectUrl := "https://" + currentHost + ":1618/" + c.Param("containerId") + "/"

	accessTokenCookie := new(http.Cookie)
	accessTokenCookie.Name = "os-access-token"
	accessTokenCookie.Value = accessToken.String()
	accessTokenCookie.Expires = time.Now().Add(3 * time.Hour)
	accessTokenCookie.Path = redirectUrl
	accessTokenCookie.HttpOnly = true
	accessTokenCookie.Secure = true
	accessTokenCookie.SameSite = http.SameSiteStrictMode
	c.SetCookie(accessTokenCookie)

	return c.Redirect(http.StatusFound, redirectUrl)
}

func (controller *ContainerController) parsePortBindings(
	rawPortBindings []interface{},
) (portBindings []valueObject.PortBinding) {
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

		rawPublicPort, rawPublicPortExists := rawPortBindingMap["publicPort"]
		if rawPublicPortExists {
			publicPort, err := valueObject.NewNetworkPort(rawPublicPort)
			if err != nil {
				log.Printf("[%d] %s", bindingIndex, err.Error())
				continue
			}
			if len(portBindingStr) > 0 {
				portBindingStr += ":"
			}
			portBindingStr += publicPort.String()
		}

		rawContainerPort, rawContainerPortExists := rawPortBindingMap["containerPort"]
		if rawContainerPortExists {
			containerPort, err := valueObject.NewNetworkPort(rawContainerPort)
			if err != nil {
				log.Printf("[%d] %s", bindingIndex, err.Error())
				continue
			}
			if len(portBindingStr) > 0 {
				portBindingStr += ":"
			}
			portBindingStr += containerPort.String()
		}

		rawProtocol, assertOk := rawPortBindingMap["protocol"].(string)
		if assertOk && (rawPublicPortExists || rawContainerPortExists) {
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
			if len(portBindingStr) > 0 {
				portBindingStr += ":"
			}
			portBindingStr += privatePort.String()
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
) (containerEnvs []valueObject.ContainerEnv) {
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
// @Summary      CreateContainer
// @Description  Create a new container.
// @Tags         container
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createContainerDto 	  body    dto.CreateContainer  true  "Only accountId, hostname and imageAddress are required.<br />When specifying portBindings, only 'publicPort' OR 'serviceName' is required.<br />'launchScript' must be base64 encoded (if any)."
// @Success      201 {object} object{} "ContainerCreationScheduled"
// @Router       /v1/container/ [post]
func (controller *ContainerController) Create(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	if requestBody["accountId"] == nil {
		requestBody["accountId"] = requestBody["operatorAccountId"]
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

		launchScript, err := valueObject.NewLaunchScriptFromEncodedContent(scriptEncodedContent)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}

		requestBody["launchScript"] = launchScript
	}

	return apiHelper.ServiceResponseWrapper(
		c,
		controller.containerService.Create(requestBody, true),
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

	if requestBody["accountId"] == nil {
		requestBody["accountId"] = requestBody["operatorAccountId"]
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
		"accountId":         c.Param("accountId"),
		"containerId":       c.Param("containerId"),
		"operatorAccountId": c.Get("accountId"),
		"operatorIpAddress": c.RealIP(),
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.containerService.Delete(requestBody),
	)
}
