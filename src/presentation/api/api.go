package api

import (
	apiInit "github.com/goinfinite/fleet/src/presentation/api/init"
	apiMiddleware "github.com/goinfinite/fleet/src/presentation/api/middleware"
	"github.com/goinfinite/fleet/src/presentation/shared"
	"github.com/labstack/echo/v4"
	_ "github.com/swaggo/echo-swagger/example/docs"
)

// @title			FleetApi
// @version			0.0.1
// @description		Infinite FleetManager API
// @termsOfService	https://goinfinite.net/tos/

// @contact.name	Infinite Engineering
// @contact.url		https://goinfinite.net/
// @contact.email	eng+swagger@goinfinite.net

// @license.name  INFINITE CLOUD COMPUTING © 2023. All Rights Reserved.
// @license.url   https://goinfinite.net/tos/

// @securityDefinitions.apikey	Bearer
// @in 							header
// @name						Authorization
// @description					Type "Bearer" + JWT token or API key.

// @host		localhost:3141
// @BasePath	/v1
func ApiInit() {
	shared.CheckEnvs()

	e := echo.New()

	basePath := "/v1"
	baseRoute := e.Group(basePath)

	e.Pre(apiMiddleware.TrailingSlash(basePath))
	e.Use(apiMiddleware.PanicHandler)
	e.Use(apiMiddleware.SetDefaultHeaders)

	dbSvc := apiInit.DatabaseService()
	e.Use(apiMiddleware.SetDatabaseService(dbSvc))

	e.Use(apiMiddleware.Auth(basePath))

	registerApiRoutes(baseRoute, dbSvc)

	e.Start(":3141")
}
