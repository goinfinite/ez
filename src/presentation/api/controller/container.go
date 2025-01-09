package apiController

import (
	"errors"
	"log/slog"
	"mime/multipart"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/useCase"
	"github.com/goinfinite/ez/src/domain/valueObject"
	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
	"github.com/goinfinite/ez/src/infra"
	"github.com/goinfinite/ez/src/infra/db"
	apiHelper "github.com/goinfinite/ez/src/presentation/api/helper"
	"github.com/goinfinite/ez/src/presentation/service"
	sharedHelper "github.com/goinfinite/ez/src/presentation/shared/helper"
	"github.com/labstack/echo/v4"
)

type ContainerController struct {
	containerService *service.ContainerService
	persistentDbSvc  *db.PersistentDatabaseService
	trailDbSvc       *db.TrailDatabaseService
}

func NewContainerController(
	persistentDbSvc *db.PersistentDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) *ContainerController {
	return &ContainerController{
		containerService: service.NewContainerService(
			persistentDbSvc, trailDbSvc,
		),
		persistentDbSvc: persistentDbSvc,
		trailDbSvc:      trailDbSvc,
	}
}

// ReadContainers	 godoc
// @Summary      ReadContainers
// @Description  List containers.
// @Tags         container
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        containerId query  string  false  "ContainerId"
// @Param        containerAccountId query  uint  false  "ContainerAccountId"
// @Param        containerHostname query  string  false  "ContainerHostname"
// @Param        containerStatus query  bool  false  "ContainerStatus"
// @Param        containerImageId query  string  false  "ContainerImageId"
// @Param        containerImageAddress query  string  false  "ContainerImageAddress"
// @Param        containerImageHash query  string  false  "ContainerImageHash"
// @Param        containerPortBindings query  string  false  "ContainerPortBindings"
// @Param        containerRestartPolicy query  string  false  "ContainerRestartPolicy"
// @Param        containerProfileId query  uint  false  "ContainerProfileId"
// @Param        containerEnv query  string  false  "ContainerEnv"
// @Param        createdBeforeAt query  string  false  "CreatedBeforeAt"
// @Param        createdAfterAt query  string  false  "CreatedAfterAt"
// @Param        startedBeforeAt query  string  false  "StartedBeforeAt"
// @Param        startedAfterAt query  string  false  "StartedAfterAt"
// @Param        stoppedBeforeAt query  string  false  "StoppedBeforeAt"
// @Param        stoppedAfterAt query  string  false  "StoppedAfterAt"
// @Param        withMetrics query  bool  false  "WithMetrics"
// @Param        pageNumber query  uint  false  "PageNumber (Pagination)"
// @Param        itemsPerPage query  uint  false  "ItemsPerPage (Pagination)"
// @Param        sortBy query  string  false  "SortBy (Pagination)"
// @Param        sortDirection query  string  false  "SortDirection (Pagination)"
// @Param        lastSeenId query  string  false  "LastSeenId (Pagination)"
// @Success      200 {object} dto.ReadContainersResponse
// @Router       /v1/container/ [get]
// @Router       /v1/container/metrics/ [get]
func (controller *ContainerController) Read(c echo.Context) error {
	requestBody := map[string]interface{}{}
	queryParameters := []string{
		"containerId", "containerAccountId", "containerHostname", "containerStatus",
		"containerImageId", "containerImageAddress", "containerImageHash",
		"containerPortBindings", "containerRestartPolicy", "containerProfileId",
		"containerEnv", "createdBeforeAt", "createdAfterAt",
		"startedBeforeAt", "startedAfterAt", "stoppedBeforeAt", "stoppedAfterAt",
		"pageNumber", "itemsPerPage", "sortBy", "sortDirection", "lastSeenId",
	}
	for _, paramName := range queryParameters {
		paramValue := c.QueryParam(paramName)
		if paramValue == "" {
			continue
		}

		if paramName == "containerId" {
			requestBody["containerId"] = sharedHelper.StringSliceValueObjectParser(
				requestBody["containerId"], valueObject.NewContainerId,
			)
		}

		if paramName == "containerAccountId" {
			requestBody["containerAccountId"] = sharedHelper.StringSliceValueObjectParser(
				requestBody["containerAccountId"], valueObject.NewAccountId,
			)
		}

		if paramName == "containerPortBindings" {
			requestBody[paramName] = controller.parsePortBindings(paramValue)
			continue
		}

		if paramName == "containerEnv" {
			requestBody[paramName] = controller.parseContainerEnvs(paramValue)
			continue
		}

		requestBody[paramName] = strings.Trim(paramValue, "\"")
	}

	if c.Request().URL.Path == "/v1/container/metrics/" {
		requestBody["withMetrics"] = true
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.containerService.Read(requestBody),
	)
}

// CreateContainerSessionToken	 godoc
// @Summary      CreateContainerSessionToken
// @Description  Creates a session token for the specified container and redirects to Infinite OS dashboard (if shouldRedirect is not false).
// @Tags         container
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        accountId	path	string	true	"AccountId"
// @Param        containerId	path	string	true	"ContainerId"
// @Param        shouldRedirect	query	bool	false	"ShouldRedirect (default/empty is true)"
// @Success      200 {object} valueObject.AccessTokenValue "If shouldRedirect is set to false, the updated session token is returned."
// @Success      302 {string} string "A redirect to Infinite OS dashboard (:1618/{containerId}/)."
// @Failure      500 {string} string "Container is not found, not running or isn't Infinite OS."
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

func (controller *ContainerController) transformPortBindingStringIntoMap(
	rawPortBindings string,
) []map[string]interface{} {
	portBindingMaps := []map[string]interface{}{}

	rawPortBindingsParts := strings.Split(rawPortBindings, ";")
	if len(rawPortBindingsParts) == 0 {
		rawPortBindingsParts = append(rawPortBindingsParts, rawPortBindings)
	}

	for _, rawPortBinding := range rawPortBindingsParts {
		portBindingMap := map[string]interface{}{}
		rawPortBindingParts := strings.Split(rawPortBinding, "|")

		possibleFieldKeys := []string{
			"serviceName", "publicPort", "containerPort", "protocol", "privatePort",
		}
		for partIndex, fieldValue := range rawPortBindingParts {
			if partIndex >= len(possibleFieldKeys) {
				break
			}

			portBindingMap[possibleFieldKeys[partIndex]] = fieldValue
		}

		portBindingMaps = append(portBindingMaps, portBindingMap)
	}

	return portBindingMaps
}

func (controller *ContainerController) parsePortBindingMap(
	rawPortBindingMap map[string]interface{},
) (portBindings []valueObject.PortBinding, err error) {
	portBindingStr := ""
	rawServiceName, assertOk := rawPortBindingMap["serviceName"].(string)
	if assertOk {
		// ServiceName won't be validated here cause it may actually be publicPort.
		// Don't worry as it will be validated later on.
		portBindingStr += rawServiceName
	}

	rawPublicPort, rawPublicPortExists := rawPortBindingMap["publicPort"]
	if rawPublicPortExists {
		publicPort, err := valueObject.NewNetworkPort(rawPublicPort)
		if err != nil {
			return portBindings, errors.New("InvalidPublicPort")
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
			return portBindings, errors.New("InvalidContainerPort")
		}

		if len(portBindingStr) > 0 {
			portBindingStr += ":"
		}
		portBindingStr += containerPort.String()
	}

	rawProtocol, rawProtocolExists := rawPortBindingMap["protocol"]
	if rawProtocolExists && (rawPublicPortExists || rawContainerPortExists) {
		protocol, err := valueObject.NewNetworkProtocol(rawProtocol)
		if err != nil {
			return portBindings, errors.New("InvalidProtocol")
		}

		portBindingStr += "/" + protocol.String()
	}

	rawPrivatePort, rawPrivatePortExists := rawPortBindingMap["privatePort"]
	if rawPrivatePortExists {
		privatePort, err := valueObject.NewNetworkPort(rawPrivatePort)
		if err != nil {
			return portBindings, errors.New("InvalidPrivatePort")
		}

		if len(portBindingStr) > 0 {
			portBindingStr += ":"
		}
		portBindingStr += privatePort.String()
	}

	return valueObject.NewPortBindingFromString(portBindingStr)
}

// PortBindings has multiple possible structures which this parser can handle:
// "serviceName" (string) OR "publicPort" (string, but actually an uint)
// "serviceName|publicPort" (string, pipe separated values)
// "serviceName|publicPort;serviceName|publicPort" (string slice, semicolon separated items)
// { "serviceName": "serviceName", "publicPort": "publicPort"} (map[string]interface{})
// [{ "serviceName": "serviceName", "publicPort": "publicPort"}] (map[string]interface{} slice)
// Besides the mentioned fields, it can also have "containerPort", "protocol" and "privatePort".
func (controller *ContainerController) parsePortBindings(
	rawPortBindings any,
) (portBindings []valueObject.PortBinding) {
	rawPortBindingsSlice := []interface{}{}

	switch rawPortBindingsValue := rawPortBindings.(type) {
	case map[string]interface{}:
		rawPortBindingsSlice = []interface{}{rawPortBindings}
	case string:
		portBindingsMaps := controller.transformPortBindingStringIntoMap(rawPortBindingsValue)
		for _, portBindingMap := range portBindingsMaps {
			rawPortBindingsSlice = append(rawPortBindingsSlice, portBindingMap)
		}
	case []interface{}:
		rawPortBindingsSlice = rawPortBindingsValue
	}

	for _, rawPortBinding := range rawPortBindingsSlice {
		rawPortBindingMap, assertOk := rawPortBinding.(map[string]interface{})
		if !assertOk {
			slog.Debug(
				"InvalidPortBindingStructure", slog.Any("rawPortBinding", rawPortBinding),
			)
			continue
		}

		portBinding, err := controller.parsePortBindingMap(rawPortBindingMap)
		if err != nil {
			slog.Debug(err.Error(), slog.Any("rawPortBinding", rawPortBinding))
		}

		portBindings = append(portBindings, portBinding...)
	}

	return portBindings
}

// Envs may come in the following structures:
// "key=value" OR "key|value" (string) OR "key=value;key=value" OR "key|value;key|value" (string, semicolon separated items)
// ["key=value", "key=value"] (string slice)
func (controller *ContainerController) parseContainerEnvs(
	envs any,
) (containerEnvs []valueObject.ContainerEnv) {
	rawEnvsSlice := []string{}

	switch envsValue := envs.(type) {
	case string:
		rawEnvsSlice = strings.Split(envsValue, ";")
	case []interface{}:
		for _, env := range envsValue {
			rawEnv, assertOk := env.(string)
			if !assertOk {
				slog.Debug("InvalidEnvStructure", slog.Any("env", env))
				continue
			}

			rawEnvsSlice = append(rawEnvsSlice, rawEnv)
		}
	}

	for _, rawEnv := range rawEnvsSlice {
		rawEnv = strings.ReplaceAll(rawEnv, "|", "=")

		containerEnv, err := valueObject.NewContainerEnv(rawEnv)
		if err != nil {
			slog.Debug(err.Error(), slog.Any("rawEnv", rawEnv))
			continue
		}

		containerEnvs = append(containerEnvs, containerEnv)
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
// @Param        archiveImageFile	formData	file	false	"ArchiveImageFile (For importing container image file, if any.)"
// @Success      201 {object} object{} "ContainerCreationScheduled"
// @Router       /v1/container/ [post]
func (controller *ContainerController) Create(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	if _, exists := requestBody["imageAddress"]; !exists {
		possibleKeys := []string{"imgAddr", "imageAddr", "imgAddress"}
		for _, possibleKey := range possibleKeys {
			if _, exists = requestBody[possibleKey]; !exists {
				continue
			}

			requestBody["imageAddress"] = requestBody[possibleKey]
		}
	}

	if requestBody["portBindings"] != nil {
		requestBody["portBindings"] = controller.parsePortBindings(
			requestBody["portBindings"],
		)
	}

	if requestBody["envs"] != nil {
		requestBody["envs"] = controller.parseContainerEnvs(requestBody["envs"])
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

	if requestBody["archiveImageFile"] != nil {
		archiveImageFiles, assertOk := requestBody["archiveImageFile"].([]*multipart.FileHeader)
		if !assertOk {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, "InvalidArchiveImageFile")
		}
		if len(archiveImageFiles) == 0 {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, "EmptyArchiveImageFile")
		}
		archiveImageFile := archiveImageFiles[0]

		operatorAccountId, err := valueObject.NewAccountId(requestBody["operatorAccountId"])
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, "InvalidOperatorAccountId")
		}

		operatorIpAddress, err := valueObject.NewIpAddress(requestBody["operatorIpAddress"])
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}

		accountId, err := valueObject.NewAccountId(requestBody["accountId"])
		if err != nil {
			requestBody["accountId"] = operatorAccountId.Uint64()
			accountId = operatorAccountId
		}

		importDto := dto.NewImportContainerImageArchiveFile(
			accountId, archiveImageFile, operatorAccountId, operatorIpAddress,
		)

		containerImageCmdRepo := infra.NewContainerImageCmdRepo(controller.persistentDbSvc)
		accountQueryRepo := infra.NewAccountQueryRepo(controller.persistentDbSvc)
		activityRecordCmdRepo := infra.NewActivityRecordCmdRepo(controller.trailDbSvc)

		importedImageId, err := useCase.ImportContainerImageArchiveFile(
			containerImageCmdRepo, accountQueryRepo, activityRecordCmdRepo, importDto,
		)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
		}

		requestBody["imageId"] = importedImageId.String()
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

	if requestBody["containerId"] == nil {
		if _, exists := requestBody["id"]; exists {
			requestBody["containerId"] = requestBody["id"]
		}
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
