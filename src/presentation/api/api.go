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
func ApiInit(
	persistentDbSvc *db.PersistentDatabaseService,
	transientDbSvc *db.TransientDatabaseService,
) {
	sharedMiddleware.CheckEnvs()

	e := echo.New()

	e.Pre(apiMiddleware.TrailingSlash())
	e.Use(apiMiddleware.PanicHandler)
	e.Use(apiMiddleware.SetDefaultHeaders)
	e.Use(apiMiddleware.SetDatabaseServices(persistentDbSvc, transientDbSvc))

	sharedMiddleware.InvalidLicenseBlocker(persistentDbSvc, transientDbSvc)

	apiInit.BootContainers(persistentDbSvc)

	e.Use(apiMiddleware.Auth())

	registerApiRoutes(e, persistentDbSvc, transientDbSvc)

	e.Start(":3141")
}
