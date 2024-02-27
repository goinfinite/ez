package api

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/infra/db"
	apiHelper "github.com/speedianet/control/src/presentation/api/helper"
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

// @BasePath	/_/api
func ApiInit(
	persistentDbSvc *db.PersistentDatabaseService,
	transientDbSvc *db.TransientDatabaseService,
) {
	sharedMiddleware.CheckEnvs()

	e := echo.New()
	apiBasePath := "/_/api"

	e.Pre(apiMiddleware.AddTrailingSlash(apiBasePath))
	e.Use(apiMiddleware.PanicHandler)
	e.Use(apiMiddleware.SetDefaultHeaders(apiBasePath))
	e.Use(apiMiddleware.SetDatabaseServices(persistentDbSvc, transientDbSvc))

	sharedMiddleware.InvalidLicenseBlocker(persistentDbSvc, transientDbSvc)

	apiInit.BootContainers(persistentDbSvc)

	e.Use(apiMiddleware.Auth(apiBasePath))

	registerApiRoutes(e, persistentDbSvc, transientDbSvc)

	httpServer := http.Server{
		Addr:     ":3141",
		Handler:  e,
		ErrorLog: apiHelper.NewCustomLogger(),
	}

	pkiDir := "/var/speedia/pki"
	certFile := pkiDir + "/control.crt"
	keyFile := pkiDir + "/control.key"

	err := httpServer.ListenAndServeTLS(certFile, keyFile)
	if err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
