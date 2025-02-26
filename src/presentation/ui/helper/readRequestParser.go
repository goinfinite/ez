package uiHelper

import (
	"reflect"
	"strings"

	"github.com/labstack/echo/v4"
)

func ReadRequestParser(
	echoContext echo.Context,
	itemName string,
	readRequestDto interface{},
) map[string]interface{} {
	requestParamsMap := map[string]interface{}{}

	if readRequestDto == nil {
		return requestParamsMap
	}

	structType := reflect.TypeOf(readRequestDto)
	for fieldIndex := range structType.NumField() {
		structField := structType.Field(fieldIndex)
		if structField.Type.Kind() == reflect.Slice {
			continue
		}

		requestParam := echoContext.QueryParam(itemName + structField.Name)
		if requestParam != "" {
			fieldNameCamelCase := strings.ToLower(structField.Name[:1]) + structField.Name[1:]
			requestParamsMap[fieldNameCamelCase] = requestParam
		}
	}

	return requestParamsMap
}
