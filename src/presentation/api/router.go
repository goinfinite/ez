package api

import (
	_ "embed"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/infra/db"
	apiController "github.com/speedianet/control/src/presentation/api/controller"
	"github.com/speedianet/control/src/presentation/ui"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/speedianet/control/src/presentation/api/docs"
)

func swaggerRoute(v1BaseRoute *echo.Group) {
	swaggerGroup := v1BaseRoute.Group("/swagger")
	swaggerGroup.GET("/*", echoSwagger.WrapHandler)
}

func authRoutes(v1BaseRoute *echo.Group) {
	authGroup := v1BaseRoute.Group("/auth")
	authGroup.POST("/login/", apiController.AuthLoginController)
}

func accountRoutes(v1BaseRoute *echo.Group, persistentDbSvc *db.PersistentDatabaseService) {
	accountGroup := v1BaseRoute.Group("/account")
	accountGroup.GET("/", apiController.GetAccountsController)
	accountGroup.POST("/", apiController.AddAccountController)
	accountGroup.PUT("/", apiController.UpdateAccountController)
	accountGroup.DELETE("/:accountId/", apiController.DeleteAccountController)
	go apiController.AutoUpdateAccountsQuotaUsageController(persistentDbSvc)
}

func containerRoutes(v1BaseRoute *echo.Group) {
	containerGroup := v1BaseRoute.Group("/container")
	containerGroup.GET("/", apiController.GetContainersController)
	containerGroup.GET("/metrics/", apiController.GetContainersWithMetricsController)
	containerGroup.POST("/", apiController.AddContainerController)
	containerGroup.PUT("/", apiController.UpdateContainerController)
	containerGroup.DELETE(
		"/:accountId/:containerId/",
		apiController.DeleteContainerController,
	)

	containerProfileGroup := containerGroup.Group("/profile")
	containerProfileGroup.GET("/", apiController.GetContainerProfilesController)
	containerProfileGroup.POST("/", apiController.AddContainerProfileController)
	containerProfileGroup.PUT("/", apiController.UpdateContainerProfileController)
	containerProfileGroup.DELETE(
		"/:profileId/",
		apiController.DeleteContainerProfileController,
	)

	containerRegistryGroup := containerGroup.Group("/registry")
	containerRegistryGroup.GET("/image/", apiController.GetContainerRegistryImagesController)
	containerRegistryGroup.GET(
		"/image/tagged/",
		apiController.GetContainerRegistryTaggedImageController,
	)
}

func licenseRoutes(
	v1BaseRoute *echo.Group,
	persistentDbSvc *db.PersistentDatabaseService,
	transientDbSvc *db.TransientDatabaseService,
) {
	licenseGroup := v1BaseRoute.Group("/license")
	licenseGroup.GET("/", apiController.GetLicenseInfoController)
	go apiController.AutoLicenseValidationController(persistentDbSvc, transientDbSvc)
}

func mappingRoutes(v1BaseRoute *echo.Group) {
	mappingGroup := v1BaseRoute.Group("/mapping")
	mappingGroup.GET("/", apiController.GetMappingsController)
	mappingGroup.POST("/", apiController.AddMappingController)
	mappingGroup.DELETE("/:mappingId/", apiController.DeleteMappingController)
	mappingGroup.POST("/target/", apiController.AddMappingTargetController)
	mappingGroup.DELETE(
		"/:mappingId/target/:targetId/",
		apiController.DeleteMappingTargetController,
	)
}

func o11yRoutes(v1BaseRoute *echo.Group) {
	o11yGroup := v1BaseRoute.Group("/o11y")
	o11yGroup.GET("/overview/", apiController.O11yOverviewController)
}

func uiRoute(e *echo.Echo) {
	uiPath := "/_/"

	httpHandler := http.StripPrefix(uiPath, ui.UiFs())
	e.GET(uiPath+"*", echo.WrapHandler(httpHandler))
}

func registerApiRoutes(
	e *echo.Echo,
	persistentDbSvc *db.PersistentDatabaseService,
	transientDbSvc *db.TransientDatabaseService,
) {
	uiRoute(e)
	v1BaseRoute := e.Group("/api/v1")

	swaggerRoute(v1BaseRoute)
	authRoutes(v1BaseRoute)
	accountRoutes(v1BaseRoute, persistentDbSvc)
	containerRoutes(v1BaseRoute)
	licenseRoutes(v1BaseRoute, persistentDbSvc, transientDbSvc)
	mappingRoutes(v1BaseRoute)
	o11yRoutes(v1BaseRoute)
}
