package api

import (
	"github.com/labstack/echo/v4"
	"github.com/speedianet/sfm/src/domain/valueObject"
	"github.com/speedianet/sfm/src/infra"
	apiMiddleware "github.com/speedianet/sfm/src/presentation/api/middleware"
	"github.com/speedianet/sfm/src/presentation/shared"
	_ "github.com/swaggo/echo-swagger/example/docs"
)

// @title			SfmApi
// @version			0.0.1
// @description		Speedia FleetManager API
// @termsOfService	https://speedia.net/tos/

// @contact.name	Speedia Engineering
// @contact.url		https://speedia.net/
// @contact.email	eng+swagger@speedia.net

// @license.name  SPEEDIA WEB SERVICES © 2023. All Rights Reserved.
// @license.url   https://speedia.net/tos/

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
