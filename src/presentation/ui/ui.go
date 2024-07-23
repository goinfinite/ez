package ui

import (
	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/infra/db"
)

func UiInit(
	e *echo.Echo,
	persistentDbSvc *db.PersistentDatabaseService,
	transientDbSvc *db.TransientDatabaseService,
) {
	basePath := ""
	baseRoute := e.Group(basePath)

	router := NewRouter(baseRoute, persistentDbSvc, transientDbSvc)
	router.RegisterRoutes()
}
