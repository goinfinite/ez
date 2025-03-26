package apiHelper

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func StringDotNotationToHierarchicalMap(
	hierarchicalMap map[string]interface{}, remainingKeys []string, finalValue string,
) map[string]interface{} {
	if len(remainingKeys) == 1 {
		hierarchicalMap[remainingKeys[0]] = finalValue
		return hierarchicalMap
	}

	parentKey := remainingKeys[0]
	nextKeys := remainingKeys[1:]

	if _, exists := hierarchicalMap[parentKey]; !exists {
		hierarchicalMap[parentKey] = make(map[string]interface{})
	}

	hierarchicalMap[parentKey] = StringDotNotationToHierarchicalMap(
		hierarchicalMap[parentKey].(map[string]interface{}), nextKeys, finalValue,
	)

	return hierarchicalMap
}

func parseFormKeyValues(formData map[string][]string) map[string]interface{} {
	formKeyValues := map[string]interface{}{}

	for formKey, keyValues := range formData {
		keyValue := ""
		for _, value := range keyValues {
			if value == "" {
				continue
			}

			keyValue = value
			break
		}
		if keyValue == "" {
			continue
		}

		isNestedKey := strings.Contains(formKey, ".")
		if !isNestedKey {
			formKeyValues[formKey] = keyValue
			continue
		}

		keyParts := strings.Split(formKey, ".")
		if len(keyParts) < 2 {
			continue
		}

		formKeyValues = StringDotNotationToHierarchicalMap(formKeyValues, keyParts, keyValue)
	}

	return formKeyValues
}

func ReadRequestBody(c echo.Context) (map[string]interface{}, error) {
	requestBody := map[string]interface{}{}

	contentType := c.Request().Header.Get("Content-Type")

	switch {
	case strings.HasPrefix(contentType, "application/json"):
		if err := c.Bind(&requestBody); err != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "InvalidJsonBody")
		}
	case strings.HasPrefix(contentType, "application/x-www-form-urlencoded"):
		formData, err := c.FormParams()
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "InvalidFormData")
		}

		requestBody = parseFormKeyValues(formData)
	case strings.HasPrefix(contentType, "multipart/form-data"):
		multipartForm, err := c.MultipartForm()
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "InvalidMultipartFormData")
		}

		requestBody = parseFormKeyValues(multipartForm.Value)

		for fileKey, fileValue := range multipartForm.File {
			if len(fileValue) == 0 {
				continue
			}

			requestBody[fileKey] = fileValue
		}
	default:
		return nil, echo.NewHTTPError(http.StatusBadRequest, "InvalidContentType")
	}

	for queryParamName, queryParamValues := range c.QueryParams() {
		requestBody[queryParamName] = queryParamValues[0]
	}

	for _, paramName := range c.ParamNames() {
		requestBody[paramName] = c.Param(paramName)
	}

	if len(requestBody) == 0 {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "EmptyRequestBody")
	}

	requestBody["operatorAccountId"] = c.Get("accountId")
	requestBody["operatorIpAddress"] = c.RealIP()

	if requestBody["accountId"] == nil {
		requestBody["accountId"] = requestBody["operatorAccountId"]
	}

	return requestBody, nil
}
