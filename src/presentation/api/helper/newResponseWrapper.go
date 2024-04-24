package apiHelper

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/presentation/service"
)

type newFormattedResponse struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body"`
}

// TODO: Remove previous wrapper once this one is used everywhere.
func NewResponseWrapper(
	c echo.Context,
	serviceOutput service.ServiceOutput,
) error {
	responseStatus := http.StatusOK
	switch serviceOutput.Status {
	case service.Created:
		responseStatus = http.StatusCreated
	case service.MultiStatus:
		responseStatus = http.StatusMultiStatus
	case service.UserError:
		responseStatus = http.StatusBadRequest
	case service.InfraError:
		responseStatus = http.StatusInternalServerError
	}

	formattedResponse := newFormattedResponse{
		Status: responseStatus,
		Body:   serviceOutput.Body,
	}
	return c.JSON(responseStatus, formattedResponse)
}
