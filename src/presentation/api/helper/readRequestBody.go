package apiHelper

import (
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

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
		for k, v := range formData {
			if len(v) > 0 {
				requestBody[k] = v[0]
			}
		}
	case strings.HasPrefix(contentType, "multipart/form-data"):
		multipartFormData, err := c.MultipartForm()
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "InvalidMultipartFormData")
		}

		for k, v := range multipartFormData.Value {
			if len(v) > 0 {
				requestBody[k] = v[0]
			}
		}

		if len(multipartFormData.File) > 0 {
			requestBodyFiles := map[string]*multipart.FileHeader{}

			for k, v := range multipartFormData.File {
				if len(v) > 0 {
					requestBodyFiles[k] = v[0]
				}
			}

			requestBody["files"] = requestBodyFiles
		}
	default:
		return nil, echo.NewHTTPError(http.StatusBadRequest, "InvalidContentType")
	}

	if len(requestBody) == 0 {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "EmptyRequestBody")
	}

	requestBody["operatorAccountId"] = c.Get("accountId")
	requestBody["ipAddress"] = c.RealIP()

	return requestBody, nil
}
