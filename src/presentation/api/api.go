package api

import (
	"github.com/goinfinite/ez/src/infra/db"
	apiInit "github.com/goinfinite/ez/src/presentation/api/init"
	apiMiddleware "github.com/goinfinite/ez/src/presentation/api/middleware"
	"github.com/labstack/echo/v4"
)

const (
	ApiBasePath string = "/api"
)

// @title						ezApi
// @version					0.0.8
// @description			Infinite Ez API
// @termsOfService	https://goinfinite.org/tos/

// @contact.name	Infinite Engineering
// @contact.url		https://goinfinite.org/
// @contact.email	eng+swagger@goinfinite.org

// @license.name  INFINITE COMPUTACAO EM NUVEM LTDA © 2024. All Rights Reserved.
// @license.url   https://goinfinite.org/tos/

// @securityDefinitions.apikey	Bearer
// @in													header
// @name												Authorization
// @description									Type "Bearer" + JWT token or API key.

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
