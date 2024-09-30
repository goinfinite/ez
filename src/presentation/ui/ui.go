package ui

import (
	"github.com/goinfinite/ez/src/infra/db"
	uiMiddleware "github.com/goinfinite/ez/src/presentation/ui/middleware"
	"github.com/labstack/echo/v4"
)

func UiInit(
	e *echo.Echo,
	persistentDbSvc *db.PersistentDatabaseService,
	transientDbSvc *db.TransientDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) {
	basePath := ""
	baseRoute := e.Group(basePath)

	e.Use(uiMiddleware.Authentication(persistentDbSvc))

	router := NewRouter(baseRoute, persistentDbSvc, transientDbSvc, trailDbSvc)
	router.RegisterRoutes()
}
