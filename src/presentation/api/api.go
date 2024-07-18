package api

import (
	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/infra/db"
	apiInit "github.com/speedianet/control/src/presentation/api/init"
	apiMiddleware "github.com/speedianet/control/src/presentation/api/middleware"
)

const (
	ApiBasePath string = "/api"
)

// @title			ControlApi
// @version			0.0.5
// @description		Speedia Control API
// @termsOfService	https://speedia.net/tos/

// @contact.name	Speedia Engineering
// @contact.url		https://speedia.net/
// @contact.email	eng+swagger@speedia.net

// @license.name  SPEEDIA WEB SERVICES, LLC © 2024. All Rights Reserved.
// @license.url   https://speedia.net/tos/

// @securityDefinitions.apikey	Bearer
// @in 							header
// @name						Authorization
// @description					Type "Bearer" + JWT token or API key.

// @BasePath	/api
func ApiInit(
	e *echo.Echo,
	persistentDbSvc *db.PersistentDatabaseService,
	transientDbSvc *db.TransientDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) {
	baseRoute := e.Group(ApiBasePath)

	e.Pre(apiMiddleware.AddTrailingSlash(ApiBasePath))
	e.Use(apiMiddleware.PanicHandler)
	e.Use(apiMiddleware.SetDefaultHeaders(ApiBasePath))
	e.Use(apiMiddleware.SetDatabaseServices(
		persistentDbSvc, transientDbSvc, trailDbSvc,
	))
	e.Use(apiMiddleware.ReadOnlyMode(ApiBasePath))

	apiInit.BootContainers(persistentDbSvc)

	e.Use(apiMiddleware.Auth(ApiBasePath))

	router := NewRouter(baseRoute, persistentDbSvc, transientDbSvc, trailDbSvc)
	router.RegisterRoutes()
}
