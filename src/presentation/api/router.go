package api

import (
	_ "embed"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/infra/db"
	apiController "github.com/speedianet/control/src/presentation/api/controller"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/speedianet/control/src/presentation/api/docs"
)

func swaggerRoute(baseRoute *echo.Group) {
	swaggerGroup := baseRoute.Group("/swagger")
	swaggerGroup.GET("/*", echoSwagger.WrapHandler)
}

func authRoutes(baseRoute *echo.Group) {
	authGroup := baseRoute.Group("/auth")
	authGroup.POST("/login/", apiController.AuthLoginController)
}

func accountRoutes(baseRoute *echo.Group, persistentDbSvc *db.PersistentDatabaseService) {
	accountGroup := baseRoute.Group("/account")
	accountGroup.GET("/", apiController.GetAccountsController)
	accountGroup.POST("/", apiController.AddAccountController)
	accountGroup.PUT("/", apiController.UpdateAccountController)
	accountGroup.DELETE("/:accountId/", apiController.DeleteAccountController)
	go apiController.AutoUpdateAccountsQuotaUsageController(persistentDbSvc)
}

func containerRoutes(baseRoute *echo.Group) {
	containerGroup := baseRoute.Group("/container")
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
	baseRoute *echo.Group,
	persistentDbSvc *db.PersistentDatabaseService,
	transientDbSvc *db.TransientDatabaseService,
) {
	licenseGroup := baseRoute.Group("/license")
	licenseGroup.GET("/", apiController.GetLicenseInfoController)
	go apiController.AutoLicenseValidationController(persistentDbSvc, transientDbSvc)
}

func mappingRoutes(baseRoute *echo.Group) {
	mappingGroup := baseRoute.Group("/mapping")
	mappingGroup.GET("/", apiController.GetMappingsController)
	mappingGroup.POST("/", apiController.AddMappingController)
	mappingGroup.DELETE("/:mappingId/", apiController.DeleteMappingController)
	mappingGroup.POST("/target/", apiController.AddMappingTargetController)
	mappingGroup.DELETE(
		"/:mappingId/target/:targetId/",
		apiController.DeleteMappingTargetController,
	)
}

func o11yRoutes(baseRoute *echo.Group) {
	o11yGroup := baseRoute.Group("/o11y")
	o11yGroup.GET("/overview/", apiController.O11yOverviewController)
}

func registerApiRoutes(
	baseRoute *echo.Group,
	persistentDbSvc *db.PersistentDatabaseService,
	transientDbSvc *db.TransientDatabaseService,
) {
	swaggerRoute(baseRoute)
	authRoutes(baseRoute)
	accountRoutes(baseRoute, persistentDbSvc)
	containerRoutes(baseRoute)
	licenseRoutes(baseRoute, persistentDbSvc, transientDbSvc)
	mappingRoutes(baseRoute)
	o11yRoutes(baseRoute)
}
