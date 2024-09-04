package ui

import (
	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/infra/db"
	uiMiddleware "github.com/speedianet/control/src/presentation/ui/middleware"
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
