package api

import (
	"github.com/goinfinite/fleet/src/domain/valueObject"
	"github.com/goinfinite/fleet/src/infra"
	apiMiddleware "github.com/goinfinite/fleet/src/presentation/api/middleware"
	"github.com/goinfinite/fleet/src/presentation/shared"
	"github.com/labstack/echo/v4"
	_ "github.com/swaggo/echo-swagger/example/docs"
)

// @title			SfmApi
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

// @host		localhost:10001
// @BasePath	/v1
func ApiInit() {
	shared.CheckEnvs()

	e := echo.New()

	basePath := "/v1"
	baseRoute := e.Group(basePath)

	e.Pre(apiMiddleware.TrailingSlash(basePath))
	e.Use(apiMiddleware.PanicHandler)
	e.Use(apiMiddleware.SetDefaultHeaders)
	e.Use(apiMiddleware.DatabaseInit())
	e.Use(apiMiddleware.Auth(basePath))

	registerApiRoutes(baseRoute)

	e.Start(":10001")

	infra.ServerCmdRepo{}.AddServerLog(
		valueObject.NewServerLogLevelPanic("info"),
		valueObject.NewServerLogOperationPanic("sfm"),
		valueObject.NewServerLogPayloadPanic("SFM backend is up and running!"),
	)
}
