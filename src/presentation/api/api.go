package api

import (
	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/infra/db"
	apiInit "github.com/speedianet/control/src/presentation/api/init"
	apiMiddleware "github.com/speedianet/control/src/presentation/api/middleware"
	sharedMiddleware "github.com/speedianet/control/src/presentation/shared/middleware"
)

// @title			ControlApi
// @version			0.0.1
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

// @host		localhost:3141
// @BasePath	/v1
func ApiInit(persistDbSvc *db.PersistentDatabaseService) {
	sharedMiddleware.CheckEnvs()

	e := echo.New()

	basePath := "/v1"
	baseRoute := e.Group(basePath)

	e.Pre(apiMiddleware.TrailingSlash(basePath))
	e.Use(apiMiddleware.PanicHandler)
	e.Use(apiMiddleware.SetDefaultHeaders)
	e.Use(apiMiddleware.SetPersistentDatabaseService(persistDbSvc))

	sharedMiddleware.InvalidLicenseBlocker(persistDbSvc)

	apiInit.BootContainers(persistDbSvc)

	e.Use(apiMiddleware.Auth(basePath))

	registerApiRoutes(baseRoute, persistDbSvc)

	e.Start(":3141")
}
